package main

import (
    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"
    "github.com/martini-contrib/binding"
    "github.com/Lupino/collect/models"
    "mime/multipart"
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
