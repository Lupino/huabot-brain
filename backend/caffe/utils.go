package caffe

import (
    "os"
    "os/exec"
)

func run(cmdName string, args... string) (error) {
    cmd := exec.Command(cmdName, args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd.Run()
}
