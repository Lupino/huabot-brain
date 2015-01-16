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
    File *multipart.FileHeader `form:"file"`
}

type TagForm struct {
    Tag  string                `form:"tag"`
}

type DatasetForm struct {
    FileForm
    TagForm
}

func api(mart *martini.ClassicMartini) {
    mart.Post(API + "/dataset", binding.Bind(DatasetForm{}), func(form DatasetForm, r render.Render) {
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
}
