package gearman

import (
    "github.com/Lupino/huabot-brain/backend"
    "github.com/mikespook/gearman-go/client"
    "encoding/json"
    "log"
    "flag"
    "time"
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
    var wait = make(chan bool)
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
        wait <- true
    }
    _, err = c.Do(funcName, workload, client.JobNormal, jobHandler)
    if err != nil {
        log.Printf("Error: process %s fail: %s\n", funcName, err)
        return nil, err
    }

    select {
    case <- time.After(time.Second * 60):
        errResult = errors.New("TimeoutError: process " + funcName + " timeout.")
    case <- wait:
    }
    close(wait)

    if bytes.Equal(result, []byte("error")) {
        errResult = errors.New("Error: process " + funcName + " fail.")
    }
    return result, errResult
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
