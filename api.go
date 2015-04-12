package main

import (
    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"
    "github.com/martini-contrib/binding"
    "github.com/Lupino/huabot-brain/backend"
    "github.com/Lupino/huabot-brain/backend/caffe"
    "mime/multipart"
    "strconv"
    "net/http"
    "io/ioutil"
    "log"
)

const API = "/api"

type FileForm struct {
    File *multipart.FileHeader `form:"file"`
    FileId int                 `form:"file_id"`
}

type PredictForm struct {
    ImgUrl string `form:"img_url" binding:"required"`
}

type TagForm struct {
    Tag  string                `form:"tag"  binding:"required"`
}

type DataTypeForm struct {
    DataType uint              `form:"data_type"`
}

type ModelForm struct {
    ModelName string `form:"model" binding:"required"`
}

func (dt DataTypeForm) Validate(errors binding.Errors, req *http.Request) (binding.Errors) {
    if dt.DataType > 2 {
        errors = append(errors, binding.Error{
            FieldNames: []string{"data_type"},
            Classification: "TypeError",
            Message: "data_type must on [0 1 2] or nil",
        })
    }
    return errors
}

type DatasetForm struct {
    FileForm
    TagForm
    DataTypeForm
    Description string `form:"description"`
}

func api(mart *martini.ClassicMartini) {
    var engine = backend.GetEngine()
    mart.Post(API + "/datasets/?", binding.Bind(DatasetForm{}), func(form DatasetForm, r render.Render) {
        var err error
        var file *backend.File
        var tag *backend.Tag
        var dataset *backend.Dataset

        if form.FileId > 0 {
            file = &backend.File{Id: form.FileId}
            if has, _ := engine.Get(file); !has {
                r.JSON(http.StatusOK, map[string]interface{}{"err": "file not exists."})
                return
            }
        } else {
            if file, err = backend.UploadFile(form.File); err != nil {
                r.JSON(http.StatusInternalServerError, map[string]interface{}{"err": err.Error()})
                return
            }
        }

        if tag, err = backend.SaveTag(form.Tag); err != nil {
            r.JSON(http.StatusInternalServerError, map[string]interface{}{"err": err.Error()})
            return
        }

        if dataset, err = backend.SaveDataset(file, tag, form.DataType, form.Description); err != nil {
            r.JSON(http.StatusInternalServerError, map[string]interface{}{"err": err.Error()})
            return
        }

        r.JSON(http.StatusOK, map[string]*backend.Dataset{"dataset": dataset})
    })

    mart.Post(API + "/datasets/(?P<dataset_id>\\d+)/?", func(req *http.Request, params martini.Params, r render.Render) {
        datasetId, _ := strconv.Atoi(params["dataset_id"])
        var dataset = new(backend.Dataset)
        if has, err := engine.Id(datasetId).Get(dataset); err != nil {
            r.JSON(http.StatusInternalServerError, map[string]string{"err": err.Error()})
        } else if has {
            dataset.FillObject()
            dataset.Description = req.Form.Get("description")
            if _, err := engine.Id(datasetId).Update(dataset); err != nil {
                r.JSON(http.StatusInternalServerError, map[string]string{"err": err.Error()})
                return
            }
            r.JSON(http.StatusOK, map[string]*backend.Dataset{"dataset": dataset})
        } else {
            r.JSON(http.StatusNotFound, map[string]string{"err": "Dataset not exists."})
        }
    })

    mart.Get(API + "/datasets/(?P<dataset_id>\\d+)/?", func(params martini.Params, r render.Render) {
        datasetId, _ := strconv.Atoi(params["dataset_id"])
        var dataset = new(backend.Dataset)
        if has, err := engine.Id(datasetId).Get(dataset); err != nil {
            r.JSON(http.StatusInternalServerError, map[string]string{"err": err.Error()})
        } else if has {
            dataset.FillObject()
            r.JSON(http.StatusOK, map[string]*backend.Dataset{"dataset": dataset})
        } else {
            r.JSON(http.StatusNotFound, map[string]string{"err": "Dataset not exists."})
        }
    })

    mart.Delete(API + "/datasets/(?P<dataset_id>\\d+)/?", func(params martini.Params, r render.Render) {
        datasetId, _ := strconv.Atoi(params["dataset_id"])
        var dataset = new(backend.Dataset)
        if has, err := engine.Id(datasetId).Get(dataset); err != nil {
            r.JSON(http.StatusInternalServerError, map[string]string{"err": err.Error()})
        } else if has {
            if _, err := engine.Id(datasetId).Delete(dataset); err != nil {
                r.JSON(http.StatusInternalServerError, map[string]string{"err": err.Error()})
            } else {
                var sql string
                if dataset.DataType == backend.TRAIN {
                  sql = "update `tag` set `train_count` = `train_count` - 1 where `id` = ?"
                } else if dataset.DataType == backend.VAL {
                  sql = "update `tag` set `test_count` = `test_count` - 1 where `id` = ?"
                }
                engine.Exec(sql, dataset.TagId)
                r.JSON(http.StatusOK, map[string]string{})
            }
        } else {
            r.JSON(http.StatusNotFound, map[string]string{"err": "Dataset not exists."})
        }

    })

    mart.Get(API + "/datasets/?", func(req *http.Request, r render.Render) {
        var qs = req.URL.Query()
        var err error
        var max, limit int
        if max, err = strconv.Atoi(qs.Get("max")); err != nil {
            max = -1
        }

        if limit, err = strconv.Atoi(qs.Get("limit")); err != nil {
            limit = 10
        }

        if limit > 100 {
            limit = 100
        }

        var dataType = qs.Get("data_type")

        var tagName = qs.Get("tag")

        var datasets = make([]backend.Dataset, 0)
        var q = engine.Desc("id")
        if max > -1 {
            q = q.Where("id < ?", max)
        }

        if tagName != "" {
            tag := &backend.Tag{Name: tagName}
            has, _ := engine.Get(tag)
            if !has {
                r.JSON(http.StatusNotFound,
                       map[string]string{"err": "tag: " + tagName + " not found."})
                return
            }
            q = q.And("tag_id = ?", tag.Id)
        }

        if dataType == "train" {
            q = q.And("data_type = ?", backend.TRAIN)
        } else if dataType == "val" {
            q = q.And("data_type = ?", backend.VAL)
        }

        q = q.Limit(limit)
        if err = q.Find(&datasets); err != nil {
            r.JSON(http.StatusInternalServerError, map[string]string{"err": err.Error()})
        }

        for idx, dataset := range datasets {
            dataset.FillObject()
            datasets[idx] = dataset
        }

        r.JSON(http.StatusOK, map[string][]backend.Dataset{"datasets": datasets})
    })

    mart.Post(API + "/upload/?", binding.Bind(FileForm{}), func(form FileForm, r render.Render) {
        var err error
        var file *backend.File

        if file, err = backend.UploadFile(form.File); err != nil {
            r.JSON(http.StatusInternalServerError, map[string]interface{}{"err": err.Error()})
            return
        }

        r.JSON(http.StatusOK, map[string]*backend.File{"file": file})
    })

    mart.Post(API + "/tags/?", binding.Bind(TagForm{}), func(form TagForm, r render.Render) {
        var err error
        var tag *backend.Tag

        if tag, err = backend.SaveTag(form.Tag); err != nil {
            r.JSON(http.StatusInternalServerError, map[string]interface{}{"err": err.Error()})
        }

        r.JSON(http.StatusOK, map[string]*backend.Tag{"file": tag})
    })

    mart.Post(API + "/tags/(?P<tag_id>\\d+)/?", binding.Bind(TagForm{}), func(form TagForm, params martini.Params, r render.Render) {
        tagId, _ := strconv.Atoi(params["tag_id"])
        var tag = new(backend.Tag)
        if has, err := engine.Id(tagId).Get(tag); err != nil {
            r.JSON(http.StatusInternalServerError, map[string]string{"err": err.Error()})
        } else if has {
            if tag.Name != form.Tag {
                tag.Name = form.Tag
                if _, err := engine.Id(tagId).Update(tag); err != nil {
                    r.JSON(http.StatusInternalServerError, map[string]string{"err": err.Error()})
                    return
                }
            }
            r.JSON(http.StatusOK, map[string]*backend.Tag{"tag": tag})
        } else {
            r.JSON(http.StatusNotFound, map[string]string{"err": "Tag not exists."})
        }
    })

    mart.Get(API + "/tags/(?P<tag_id>\\d+)/?", func(params martini.Params, r render.Render) {
        tagId, _ := strconv.Atoi(params["tag_id"])
        var tag = new(backend.Tag)
        if has, err := engine.Id(tagId).Get(tag); err != nil {
            r.JSON(http.StatusInternalServerError, map[string]string{"err": err.Error()})
        } else if has {
            r.JSON(http.StatusOK, map[string]*backend.Tag{"tag": tag})
        } else {
            r.JSON(http.StatusNotFound, map[string]string{"err": "Tag not exists."})
        }
    })

    mart.Delete(API + "/tags/(?P<tag_id>\\d+)/?", func(params martini.Params, r render.Render) {
        tagId, _ := strconv.Atoi(params["tag_id"])
        var tag = new(backend.Tag)
        if has, err := engine.Id(tagId).Get(tag); err != nil {
            r.JSON(http.StatusInternalServerError, map[string]string{"err": err.Error()})
        } else if has {
            backend.DeleteTag(tag.Id)
            r.JSON(http.StatusOK, map[string]*backend.Tag{"tag": tag})
        } else {
            r.JSON(http.StatusNotFound, map[string]string{"err": "Tag not exists."})
        }
    })

    mart.Get(API + "/tags/hint/?", func(req *http.Request, r render.Render) {
        var qs = req.URL.Query()
        var word = qs.Get("word")
        var q = engine.Desc("id")
        q = q.And("name like \"%" + word + "%\"")
        q = q.Limit(5)
        var tags = make([]backend.Tag, 0)
        var err = q.Find(&tags)
        log.Printf("err: %s\n", err)
        r.JSON(http.StatusOK, map[string][]backend.Tag{"tags": tags})
    })

    mart.Get(API + "/tags/?", func(req *http.Request, r render.Render) {
        var qs = req.URL.Query()
        var err error
        var max, limit int
        if max, err = strconv.Atoi(qs.Get("max")); err != nil {
            max = -1
        }

        if limit, err = strconv.Atoi(qs.Get("limit")); err != nil {
            limit = 10
        }

        if limit > 100 {
            limit = 100
        }

        var tags = make([]backend.Tag, 0)
        var q = engine.Desc("id")
        if max > -1 {
            q = q.Where("id < ?", max)
        }
        q = q.Limit(limit)
        if err = q.Find(&tags); err != nil {
            r.JSON(http.StatusInternalServerError, map[string]string{"err": err.Error()})
        }

        r.JSON(http.StatusOK, map[string][]backend.Tag{"tags": tags})
    })

    mart.Post(API + "/solve/?", func(r render.Render) {
        caffe.Solve()
        r.JSON(http.StatusOK, caffe.LastStatus())
        return
    })

    mart.Get(API + "/solve/?", func(r render.Render) {
        r.JSON(http.StatusOK, caffe.LastStatus())
        return
    })

    mart.Delete(API + "/solve/?", func(r render.Render) {
        caffe.StopSolve()
        r.JSON(http.StatusOK, caffe.LastStatus())
        return
    })

    mart.Get(API + "/train.txt", func(r render.Render) {
        text, err := backend.ExportDataset(backend.TRAIN)
        if err != nil {
            r.Data(http.StatusInternalServerError, nil)
            return
        }
        r.Data(http.StatusOK, []byte(text))
    })

    mart.Get(API + "/val.txt", func(r render.Render) {
        text, err := backend.ExportDataset(backend.VAL)
        if err != nil {
            r.Data(http.StatusInternalServerError, nil)
            return
        }
        r.Data(http.StatusOK, []byte(text))
    })

    mart.Get(API + "/loss.png", func(r render.Render) {
        result, err := caffe.Plot("loss")
        if err != nil {
            r.Redirect("/static/images/loading.png")
            return
        }
        r.Data(http.StatusOK, result)
        return
    })

    mart.Get(API + "/acc.png", func(r render.Render) {
        result, err := caffe.Plot("acc")
        if err != nil {
            r.Redirect("/static/images/loading.png")
            return
        }
        r.Data(http.StatusOK, result)
        return
    })

    mart.Post(API + "/predict/?", binding.Bind(PredictForm{}), func(form PredictForm, r render.Render) {
        result, err := caffe.PredictUrl(form.ImgUrl)
        if err != nil {
            r.JSON(http.StatusInternalServerError, map[string]string{"err": err.Error()})
            return
        }
        r.JSON(http.StatusOK, result)
        return
    })

    mart.Get(API + "/proxy/?", func(req *http.Request, r render.Render) {
        var qs = req.URL.Query()
        var url = qs.Get("url")
        var resp *http.Response
        var err error
        if resp, err = http.Get(url); err != nil {
            r.Data(http.StatusNotFound, nil)
            return
        }
        defer resp.Body.Close()
        data, _ := ioutil.ReadAll(resp.Body)
        r.Data(http.StatusOK, data)
        return
    })

    mart.Get(API + "/models", func(r render.Render) {
        modelNames, err := caffe.ListModels()
        if err != nil {
            r.JSON(http.StatusInternalServerError, map[string]string{"err": err.Error()})
            return
        }
        r.JSON(http.StatusOK, map[string][]string{"models": modelNames})
        return
    })

    mart.Get(API + "/models/current", func(r render.Render) {
        modelName, err := caffe.GetCurrentModel()
        if err != nil {
            r.JSON(http.StatusInternalServerError, map[string]string{"err": err.Error()})
            return
        }
        r.JSON(http.StatusOK, map[string]string{"model": modelName})
        return
    })

    mart.Post(API + "/models/apply", binding.Bind(ModelForm{}), func(model ModelForm, r render.Render) {
        err := caffe.ApplyModel(model.ModelName)
        if err != nil {
            r.JSON(http.StatusInternalServerError, map[string]string{"err": err.Error()})
            return
        }
        r.JSON(http.StatusOK, map[string]string{})
        return
    })

    mart.Delete(API + "/models/(?P<modelName>[^/]+.caffemodel)", func(params martini.Params, r render.Render) {
        err := caffe.RemoveModel(params["modelName"])
        if err != nil {
            r.JSON(http.StatusInternalServerError, map[string]string{"err": err.Error()})
            return
        }
        r.JSON(http.StatusOK, map[string]string{})
        return
    })
}
