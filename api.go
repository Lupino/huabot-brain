package main

import (
    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"
    "github.com/martini-contrib/binding"
    "github.com/Lupino/huabot-brain/models"
    "mime/multipart"
    "strconv"
    "net/http"
    "io/ioutil"
    "log"
)


type FileForm struct {
    File *multipart.FileHeader `form:"file" binding:"required"`
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

}

func api(mart *martini.ClassicMartini) {
    var engine = models.GetEngine()
    mart.Post(API + "/datasets/?", binding.Bind(DatasetForm{}), func(form DatasetForm, r render.Render) {
        var err error
        var file *models.File
        var tag *models.Tag
        var dataset *models.Dataset

        if file, err = uploadFile(form.File); err != nil {
            r.JSON(http.StatusInternalServerError, map[string]interface{}{"err": err.Error()})
        }

        if tag, err = saveTag(form.Tag); err != nil {
            r.JSON(http.StatusInternalServerError, map[string]interface{}{"err": err.Error()})
        }

        if dataset, err = saveDataset(file, tag, form.DataType); err != nil {
            r.JSON(http.StatusInternalServerError, map[string]interface{}{"err": err.Error()})
        }

        r.JSON(http.StatusOK, map[string]*models.Dataset{"dataset": dataset})
    })

    mart.Post(API + "/datasets/(?P<dataset_id>\\d+)/?", binding.Bind(DataTypeForm{}), func(form DataTypeForm, params martini.Params, r render.Render) {
        datasetId, _ := strconv.Atoi(params["dataset_id"])
        var dataset = new(models.Dataset)
        if has, err := engine.Id(datasetId).Get(dataset); err != nil {
            r.JSON(http.StatusInternalServerError, map[string]string{"err": err.Error()})
        } else if has {
            dataset.FillObject()
            if dataset.DataType != form.DataType {
                dataset.DataType = form.DataType
                if _, err := engine.Id(datasetId).Update(dataset); err != nil {
                    r.JSON(http.StatusInternalServerError, map[string]string{"err": err.Error()})
                    return
                }
            }
            r.JSON(http.StatusOK, map[string]*models.Dataset{"dataset": dataset})
        } else {
            r.JSON(http.StatusNotFound, map[string]string{"err": "Dataset not exists."})
        }
    })

    mart.Get(API + "/datasets/(?P<dataset_id>\\d+)/?", func(params martini.Params, r render.Render) {
        datasetId, _ := strconv.Atoi(params["dataset_id"])
        var dataset = new(models.Dataset)
        if has, err := engine.Id(datasetId).Get(dataset); err != nil {
            r.JSON(http.StatusInternalServerError, map[string]string{"err": err.Error()})
        } else if has {
            dataset.FillObject()
            r.JSON(http.StatusOK, map[string]*models.Dataset{"dataset": dataset})
        } else {
            r.JSON(http.StatusNotFound, map[string]string{"err": "Dataset not exists."})
        }
    })

    mart.Delete(API + "/datasets/(?P<dataset_id>\\d+)/?", func(params martini.Params, r render.Render) {
        datasetId, _ := strconv.Atoi(params["dataset_id"])
        var dataset = new(models.Dataset)
        if _, err := engine.Id(datasetId).Delete(dataset); err != nil {
            r.JSON(http.StatusInternalServerError, map[string]string{"err": err.Error()})
        } else {
            r.JSON(http.StatusOK, map[string]string{})
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

        var datasets = make([]models.Dataset, 0)
        var q = engine.Desc("id")
        if max > -1 {
            q = q.Where("id < ?", max)
        }

        if tagName != "" {
            tag := &models.Tag{Name: tagName}
            has, _ := engine.Get(tag)
            if !has {
                r.JSON(http.StatusNotFound,
                       map[string]string{"err": "tag: " + tagName + " not found."})
                return
            }
            q = q.And("tag_id = ?", tag.Id)
        }

        if dataType == "train" {
            q = q.And("data_type = ?", models.TRAIN)
        } else if dataType == "val" {
            q = q.And("data_type = ?", models.VAL)
        }

        q = q.Limit(limit)
        if err = q.Find(&datasets); err != nil {
            r.JSON(http.StatusInternalServerError, map[string]string{"err": err.Error()})
        }

        for idx, dataset := range datasets {
            dataset.FillObject()
            datasets[idx] = dataset
        }

        r.JSON(http.StatusOK, map[string][]models.Dataset{"datasets": datasets})
    })

    mart.Post(API + "/upload/?", binding.Bind(FileForm{}), func(form FileForm, r render.Render) {
        var err error
        var file *models.File

        if file, err = uploadFile(form.File); err != nil {
            r.JSON(http.StatusInternalServerError, map[string]interface{}{"err": err.Error()})
        }

        r.JSON(http.StatusOK, map[string]*models.File{"file": file})
    })

    mart.Post(API + "/tags/?", binding.Bind(TagForm{}), func(form TagForm, r render.Render) {
        var err error
        var tag *models.Tag

        if tag, err = saveTag(form.Tag); err != nil {
            r.JSON(http.StatusInternalServerError, map[string]interface{}{"err": err.Error()})
        }

        r.JSON(http.StatusOK, map[string]*models.Tag{"file": tag})
    })

    mart.Post(API + "/tags/(?P<tag_id>\\d+)/?", binding.Bind(TagForm{}), func(form TagForm, params martini.Params, r render.Render) {
        tagId, _ := strconv.Atoi(params["tag_id"])
        var tag = new(models.Tag)
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
            r.JSON(http.StatusOK, map[string]*models.Tag{"tag": tag})
        } else {
            r.JSON(http.StatusNotFound, map[string]string{"err": "Tag not exists."})
        }
    })

    mart.Get(API + "/tags/(?P<tag_id>\\d+)/?", func(params martini.Params, r render.Render) {
        tagId, _ := strconv.Atoi(params["tag_id"])
        var tag = new(models.Tag)
        if has, err := engine.Id(tagId).Get(tag); err != nil {
            r.JSON(http.StatusInternalServerError, map[string]string{"err": err.Error()})
        } else if has {
            r.JSON(http.StatusOK, map[string]*models.Tag{"tag": tag})
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
        var tags = make([]models.Tag, 0)
        var err = q.Find(&tags)
        log.Printf("err: %s\n", err)
        r.JSON(http.StatusOK, map[string][]models.Tag{"tags": tags})
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

        var tags = make([]models.Tag, 0)
        var q = engine.Desc("id")
        if max > -1 {
            q = q.Where("id < ?", max)
        }
        q = q.Limit(limit)
        if err = q.Find(&tags); err != nil {
            r.JSON(http.StatusInternalServerError, map[string]string{"err": err.Error()})
        }

        r.JSON(http.StatusOK, map[string][]models.Tag{"tags": tags})
    })

    mart.Post(API + "/train/?", func(r render.Render) {
        result, err := caffeTrain()
        if err != nil {
            r.JSON(http.StatusInternalServerError, map[string]string{"err": err.Error()})
            return
        }
        r.JSON(http.StatusOK, map[string]string{"msg": result})
        return
    })

    mart.Get(API + "/train/?", func(r render.Render) {
        result, err := caffeTrainStatus()
        if err != nil {
            r.JSON(http.StatusInternalServerError, map[string]string{"err": err.Error()})
            return
        }
        r.JSON(http.StatusOK, map[string]string{"msg": result})
        return
    })

    mart.Get(API + "/train.txt", func(r render.Render) {
        text, err := loadDataset(models.TRAIN)
        if err != nil {
            r.Data(http.StatusInternalServerError, nil)
            return
        }
        r.Data(http.StatusOK, []byte(text))
    })

    mart.Get(API + "/val.txt", func(r render.Render) {
        text, err := loadDataset(models.VAL)
        if err != nil {
            r.Data(http.StatusInternalServerError, nil)
            return
        }
        r.Data(http.StatusOK, []byte(text))
    })

    mart.Get(API + "/loss.png", func(r render.Render) {
        result, err := caffeTrainPlot("loss")
        if err != nil {
            r.Data(http.StatusInternalServerError, nil)
            return
        }
        r.Data(http.StatusOK, result)
        return
    })

    mart.Get(API + "/acc.png", func(r render.Render) {
        result, err := caffeTrainPlot("acc")
        if err != nil {
            r.Data(http.StatusInternalServerError, nil)
            return
        }
        r.Data(http.StatusOK, result)
        return
    })

    mart.Post(API + "/predict/?", binding.Bind(PredictForm{}), func(form PredictForm, r render.Render) {
        result, err := caffePredict(form.ImgUrl)
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
}
