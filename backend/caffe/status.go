package caffe

import (
    "strings"
    "io/ioutil"
    "github.com/Lupino/huabot-brain/config"
)

type Status struct {
    State string `json:"status"`
    Loss  string `json:"loss"`
    Acc   string `json:"acc"`
    Model string `json:"model"`
    PredictStatus string `json:"predictStatus"`
}

func LastStatus() (status Status) {
    var err error
    status = Status{}
    status.State = "Solved"
    if IsOnSolving() {
        status.State = "Solving"
    }
    if err = run(config.RES + "/last_status.sh",
                 config.LOG_DIR + "/caffe.INFO",
                 config.RES + "/status"); err != nil {
      return
    }

    status.Model, _ = GetCurrentModel()

    if IsPredictAlive() {
        status.PredictStatus = "Started"
    } else {
        status.PredictStatus = "Stoped"
    }

    var tmp []byte
    tmp, _ = ioutil.ReadFile(config.RES + "/status.acc.txt")
    status.Acc = strings.Trim(string(tmp), "\n ")
    tmp, _ = ioutil.ReadFile(config.RES + "/status.loss.txt")
    status.Loss = strings.Trim(string(tmp), "\n ")
    if status.Acc == "" {
        status.Acc = "0"
    }
    if status.Loss == "" {
        status.Loss = "0"
    }
    return
}
