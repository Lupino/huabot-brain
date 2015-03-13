package main

import (
    "github.com/Lupino/huabot-brain/models"
    "github.com/mikespook/gearman-go/client"
    "mime/multipart"
    "crypto/sha1"
    "io"
    "os"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "sync"
    "image"
    "log"
    _ "image/png"
    _ "image/jpeg"
)

type PredictTag struct {
    Id int         `json:"id,omitempty"`
    Score float64  `json:"score,omitempty"`
    Tag models.Tag `json:"tag,omitempty"`
}
type PredictResult struct {
    BetResult []PredictTag `json:"bet_result,omitempty"`
    Time  float64          `json:"time,omitempty"`
    Error string           `json:"err,omitempty"`
}

func uploadFile(realFile *multipart.FileHeader) (file *models.File, err error) {
    var source multipart.File
    if source, err = realFile.Open(); err != nil {
        return
    }
    defer source.Close()
    var img image.Config
    if img, _, err = image.DecodeConfig(source); err != nil {
        return
    }
    source.Seek(0, 0)
    hasher := sha1.New()
    io.Copy(hasher, source)
    fileKey := hex.EncodeToString(hasher.Sum(nil))
    file = &models.File{Key: fileKey}
    var engine = models.GetEngine()
    has, _ := engine.Get(file)
    if !has {
        var dst *os.File
        if dst, err = os.Create(*UPLOADPATH + fileKey); err != nil {
            return
        }
        defer dst.Close()
        source.Seek(0, 0)
        if _, err = io.Copy(dst, source); err != nil {
            return
        }

        file.Width = img.Width
        file.Height = img.Height

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

func saveDataset(file *models.File, tag *models.Tag, dataType uint, desc string) (dataset *models.Dataset, err error) {
    dataset = &models.Dataset{FileId: file.Id, TagId: tag.Id}
    var engine = models.GetEngine()
    has, _ := engine.Get(dataset)
    if !has {
        dataset.DataType = dataType
        dataset.Description = desc
        if _, err = engine.Insert(dataset); err != nil {
            return
        }
        var sql string
        if dataType == models.TRAIN {
          sql = "update `tag` set `train_count` = `train_count` + 1 where `id` = ?"
        } else if dataType == models.VAL {
          sql = "update `tag` set `test_count` = `test_count` + 1 where `id` = ?"
        }
        engine.Exec(sql, tag.Id)
    }
    dataset.File = file
    dataset.Tag = tag
    return
}

func loadDataset(dataType uint) (text string, err error) {
    var engine = models.GetEngine()
    err = engine.Where("data_type=?", dataType).Iterate(new(models.Dataset), func(i int, bean interface{}) error {
        dataset := bean.(*models.Dataset)
        dataset.FillObject()
        text = fmt.Sprintf("%s%s %d\n", text, dataset.File.Key, dataset.TagId)
        return nil
    })
    return
}

func submit(funcName string, workload []byte) ([]byte, error) {
    var mutex sync.Mutex
    var result []byte
    var errResult error
    c, err := client.New("tcp4", *GEARMAND)
    if err != nil {
        return nil, err
    }
    defer c.Close()
    c.ErrorHandler = func(e error) {
        log.Println(e)
    }
    jobHandler := func(resp *client.Response) {
        if resp.DataType == client.WorkComplate {
            result, errResult = resp.Result()
            mutex.Unlock()
        }
    }
    _, err = c.Do(funcName, workload, client.JobNormal, jobHandler)
    if err != nil {
        log.Printf("gearman Do %s Error: %s\n", funcName, err)
        return nil, err
    }
    mutex.Lock()
    mutex.Lock()
    return result, errResult
}

func caffeTrain() (string, error) {
    result, err := submit("CAFFE:TRAIN", nil)
    if err != nil {
        return "", err
    }
    return string(result), nil
}

func caffeTrainStop() (string, error) {
    result, err := submit("CAFFE:TRAIN:STOP", nil)
    if err != nil {
        return "", err
    }
    return string(result), nil
}

func caffeTrainStatus() ([]byte, error) {
    result, err := submit("CAFFE:TRAIN:STATUS", nil)
    if err != nil {
        return nil, err
    }
    return result, nil
}

func caffeTrainPlot(plotType string) ([]byte, error) {
    return submit("CAFFE:TRAIN:PLOT", []byte(plotType))
}

func caffePredict(url string) (result PredictResult, err error) {
    var data []byte
    if data, err = submit("CAFFE:PREDICT:URL", []byte(url)); err != nil {
        return
    }
    if err = json.Unmarshal(data, &result); err != nil {
        return
    }

    var engine = models.GetEngine()
    for i, ptag := range result.BetResult {
        ptag.Tag.Id = ptag.Id
        engine.Get(&ptag.Tag)
        result.BetResult[i].Tag = ptag.Tag
    }
    return
}
