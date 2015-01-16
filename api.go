package main

import (
    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"
    "github.com/martini-contrib/binding"
    "github.com/Lupino/collect/models"
    "mime/multipart"
    "strconv"
    "net/http"
)


type FileForm struct {
    File *multipart.FileHeader `form:"file" binding:"required"`
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

        var datasets = make([]models.Dataset, 0)
        var q = engine.Desc("id")
        if max > -1 {
            q = q.Where("id < ?", max)
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
}
