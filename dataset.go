package main

import (
    "github.com/Lupino/collect/models"
    "github.com/Lupino/collect/caffe"
    "mime/multipart"
    "crypto/sha1"
    "io"
    "os"
    "encoding/hex"
    "fmt"
    "sync"
)

var trainLocker = new(sync.Mutex)
var onTraining = false

func uploadFile(realFile *multipart.FileHeader) (file *models.File, err error) {
    var source multipart.File
    if source, err = realFile.Open(); err != nil {
        return
    }
    defer source.Close()
    hasher := sha1.New()
    io.Copy(hasher, source)
    fileKey := hex.EncodeToString(hasher.Sum(nil))
    file = &models.File{Key: fileKey}
    var engine = models.GetEngine()
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

func saveTag(realTag string) (tag *models.Tag, err error) {
    tag = &models.Tag{Name: realTag}
    var engine = models.GetEngine()
    has, _ := engine.Get(tag)
    if !has {
        if _, err = engine.Insert(tag); err != nil {
            return
        }
    }
    return
}

func saveDataset(file *models.File, tag *models.Tag, dataType uint) (dataset *models.Dataset, err error) {
    dataset = &models.Dataset{FileId: file.Id, TagId: tag.Id}
    var engine = models.GetEngine()
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

func exportDataset() (err error) {

    var trainFile, valFile *os.File
    if trainFile, err = os.Create(TRAIN_FILE); err != nil {
        return
    }
    defer trainFile.Close()

    exportToFile(trainFile, models.TRAIN)

    if valFile, err = os.Create(VAL_FILE); err != nil {
        return
    }
    defer valFile.Close()

    exportToFile(valFile, models.VAL)

    return err
}

func exportToFile(file *os.File, dataType uint) (err error) {
    var engine = models.GetEngine()
    err = engine.Where("data_type=?", dataType).Iterate(new(models.Dataset), func(i int, bean interface{}) error {
        dataset := bean.(*models.Dataset)
        dataset.FillObject()
        fmt.Fprintf(file, "%s %d\n", dataset.File.Key, dataset.TagId)
        return nil
    })
    return
}

func caffeTrain() (err error) {
    if onTraining {
        return
    }

    onTraining = true
    trainLocker.Lock()
    defer (func() {
        onTraining = false
        trainLocker.Unlock()
    })()

    if err = exportDataset(); err != nil {
        return
    }

    if err = caffe.ConvertImageset("--resize_height=256", "--shuffle", "--resize_width=256", UPLOADPATH, TRAIN_FILE, TRAIN_LMDB); err != nil {
        return
    }

    if err = caffe.ConvertImageset("--resize_height=256", "--resize_width=256", "--shuffle", UPLOADPATH, VAL_FILE, VAL_LMDB); err != nil {
        return
    }

    if err = caffe.ComputeImageMean(TRAIN_LMDB, MEAN_FILE); err != nil {
        return
    }

    return caffe.Train(SOLVER_FILE)
}
