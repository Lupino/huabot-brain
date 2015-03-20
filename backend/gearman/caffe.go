package gearman

import (
    "github.com/Lupino/huabot-brain/backend"
    "github.com/mikespook/gearman-go/client"
    "encoding/json"
    "sync"
    "log"
    "flag"
    "bytes"
    "errors"
)

var GEARMAND = flag.String("gearmand", "127.0.0.1:4730", "The Gearmand server.")

type PredictTag struct {
    Id int         `json:"id,omitempty"`
    Score float64  `json:"score,omitempty"`
    Tag backend.Tag `json:"tag,omitempty"`
}
type PredictResult struct {
    BetResult []PredictTag `json:"bet_result,omitempty"`
    Time  float64          `json:"time,omitempty"`
    Error string           `json:"err,omitempty"`
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
        result, errResult = resp.Result()
        mutex.Unlock()
    }
    _, err = c.Do(funcName, workload, client.JobNormal, jobHandler)
    if err != nil {
        log.Printf("Gearman Doing %s Error: %s\n", funcName, err)
        return nil, err
    }
    mutex.Lock()
    mutex.Lock()
    if bytes.Equal(result, []byte("error")) {
        errResult = errors.New("Gearman doing " + funcName + " Error")
    }
    return result, errResult
}

func Train() (string, error) {
    result, err := submit("CAFFE:TRAIN", nil)
    if err != nil {
        return "", err
    }
    return string(result), nil
}

func Stop() (string, error) {
    result, err := submit("CAFFE:TRAIN:STOP", nil)
    if err != nil {
        return "", err
    }
    return string(result), nil
}

func Status() ([]byte, error) {
    result, err := submit("CAFFE:TRAIN:STATUS", nil)
    if err != nil {
        return nil, err
    }
    return result, nil
}

func Plot(plotType string) ([]byte, error) {
    return submit("CAFFE:TRAIN:PLOT", []byte(plotType))
}

func Predict(url string) (result PredictResult, err error) {
    var data []byte
    if data, err = submit("CAFFE:PREDICT:URL", []byte(url)); err != nil {
        return
    }
    if err = json.Unmarshal(data, &result); err != nil {
        return
    }

    var engine = backend.GetEngine()
    for i, ptag := range result.BetResult {
        ptag.Tag.Id = ptag.Id
        engine.Get(&ptag.Tag)
        result.BetResult[i].Tag = ptag.Tag
    }
    return
}
