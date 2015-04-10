package caffe

import (
    "os"
    "strings"
    "os/exec"
    "io/ioutil"
)

type Status struct {
    State string `json:"status"`
    Loss  string `json:"loss"`
    Acc   string `json:"acc"`
}

func LastStatus() (status Status) {
    var err error
    status = Status{}
    status.State = "Solved"
    if IsOnSolving() {
        status.State = "Solving"
    }
    if err = run(resoursesPath + "/last_status.sh", LOG_DIR + "/caffe.INFO", resoursesPath + "/status"); err != nil {
      return
    }


    var tmp []byte
    tmp, _ = ioutil.ReadFile(resoursesPath + "/status.acc.txt")
    status.Acc = strings.Trim(string(tmp), "\n ")
    tmp, _ = ioutil.ReadFile(resoursesPath + "/status.loss.txt")
    status.Loss = strings.Trim(string(tmp), "\n ")
    if status.Acc == "" {
        status.Acc = "0"
    }
    if status.Loss == "" {
        status.Loss = "0"
    }
    return
}

func run(cmdName string, args... string) (error) {
    cmd := exec.Command(cmdName, args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd.Run()
}
