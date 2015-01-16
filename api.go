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

type DatasetForm struct {
    FileForm
    TagForm
}

func api(mart *martini.ClassicMartini) {
    mart.Post(API + "/datasets", binding.Bind(DatasetForm{}), func(form DatasetForm, r render.Render) {
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

        if dataset, err = saveDataset(file, tag, 0); err != nil {
            r.JSON(http.StatusInternalServerError, map[string]interface{}{"err": err.Error()})
        }

        r.JSON(http.StatusOK, map[string]*models.Dataset{"dataset": dataset})
    })

    mart.Post(API + "/upload", binding.Bind(FileForm{}), func(form FileForm, r render.Render) {
        var err error
        var file *models.File

        if file, err = uploadFile(form.File); err != nil {
            r.JSON(http.StatusInternalServerError, map[string]interface{}{"err": err.Error()})
        }

        r.JSON(http.StatusOK, map[string]*models.File{"file": file})
    })

    mart.Post(API + "/tags", binding.Bind(TagForm{}), func(form TagForm, r render.Render) {
        var err error
        var tag *models.Tag

        if tag, err = saveTag(form.Tag); err != nil {
            r.JSON(http.StatusInternalServerError, map[string]interface{}{"err": err.Error()})
        }

        r.JSON(http.StatusOK, map[string]*models.Tag{"file": tag})
    })
}
