package main

import (
    "github.com/Lupino/collect/models"
    "github.com/go-xorm/xorm"
    "mime/multipart"
    "crypto/sha1"
    "io"
    "os"
    "encoding/hex"
)

func uploadFile(realFile *multipart.FileHeader, engine *xorm.Engine) (file *models.File, err error) {
    var source multipart.File
    if source, err = realFile.Open(); err != nil {
        return
    }
    defer source.Close()
    hasher := sha1.New()
    io.Copy(hasher, source)
    fileKey := hex.EncodeToString(hasher.Sum(nil))
    file = &models.File{Key: fileKey}
    has, _ := engine.Get(file)
    if !has {
        var dst *os.File
        if dst, err = os.Create(UPLOADPATH + fileKey); err != nil {
            return
        }
        defer dst.Close()
        source.Seek(0, 0)
        if _, err = io.Copy(dst, source); err != nil {
            return
        }

        if _, err = engine.Insert(file); err != nil {
            return
        }
    }
    return
}

func saveTag(realTag string, engine *xorm.Engine) (tag *models.Tag, err error) {
    tag = &models.Tag{Name: realTag}
    has, _ := engine.Get(tag)
    if !has {
        if _, err = engine.Insert(tag); err != nil {
            return
        }
    }
    return
}

func saveDataset(file *models.File, tag *models.Tag, dataType int, engine *xorm.Engine) (dataset *models.Dataset, err error) {
    dataset = &models.Dataset{FileId: file.Id, TagId: tag.Id}
    has, _ := engine.Get(dataset)
    if !has {
        dataset.DataType = dataType
        if _, err = engine.Insert(dataset); err != nil {
            return
        }
    }
    dataset.File = file
    dataset.Tag = tag
    return
}
