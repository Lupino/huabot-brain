package caffe

import (
    "os"
    "fmt"
    "bytes"
    "os/exec"
    "net/url"
    "net/http"
    "encoding/json"
    "github.com/Lupino/huabot-brain/backend"
    "github.com/Lupino/huabot-brain/config"
)

type PredictTag struct {
    Id int          `json:"id,omitempty"`
    Score float64   `json:"score,omitempty"`
    Tag backend.Tag `json:"tag,omitempty"`
}

type PredictResult struct {
    BetResult []PredictTag `json:"bet_result,omitempty"`
    Time  float64          `json:"time,omitempty"`
    Error string           `json:"err,omitempty"`
}

var predictCmd *exec.Cmd

func PredictUrl(imgUrl string) (result PredictResult, err error) {
    resp, err := http.PostForm(config.PREDICT_HOST + "/api/predict/url",
                               url.Values{"img_url": {imgUrl}})

    retBody := &bytes.Buffer{}
    if _, err = retBody.ReadFrom(resp.Body); err != nil {
        return
    }
    resp.Body.Close()
    if err = json.Unmarshal(retBody.Bytes(), &result); err != nil {
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

func StartPredict() (error) {
    if predictCmd != nil && !predictCmd.ProcessState.Exited() {
        return fmt.Errorf("Predict is alreadly running.")
    }
    predictCmd = exec.Command(config.PREDICT, "--host", config.PREDICT_HOST)
    predictCmd.Stdout = os.Stdout
    predictCmd.Stderr = os.Stderr
    return predictCmd.Start()
}

func StopPredict() (error) {
    if predictCmd != nil && predictCmd.Process != nil {
        return predictCmd.Process.Kill()
    }
    return nil
}
