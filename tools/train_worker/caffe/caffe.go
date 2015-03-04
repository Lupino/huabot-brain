package caffe

import (
    "os"
    "sync"
    "os/exec"
)

var Locker = new(sync.Mutex)

var CAFFE = "caffe"
var COMPUTER_IMAGE_MEAN = "compute_image_mean"
var CONVERT_IMAGESET = "convert_imageset"

var currentCmd *exec.Cmd

func Run(args... string) error {
    Locker.Lock()
    defer Locker.Unlock()
    currentCmd = exec.Command(CAFFE, args...)
    currentCmd.Stdout = os.Stdout
    currentCmd.Stderr = os.Stderr
    return currentCmd.Run()
}

func ComputeImageMean(sourceFile, binaryFile string) error {
    Locker.Lock()
    defer Locker.Unlock()
    currentCmd = exec.Command(COMPUTER_IMAGE_MEAN, sourceFile, binaryFile)
    currentCmd.Stdout = os.Stdout
    currentCmd.Stderr = os.Stderr
    return currentCmd.Run()
}

func ConvertImageset(args... string) error {
    Locker.Lock()
    defer Locker.Unlock()
    currentCmd = exec.Command(CONVERT_IMAGESET, args...)
    currentCmd.Stdout = os.Stdout
    currentCmd.Stderr = os.Stderr
    return currentCmd.Run()
}

func Kill() (err error) {
    if currentCmd != nil && currentCmd.Process != nil {
        return currentCmd.Process.Kill()
    }
    return
}
