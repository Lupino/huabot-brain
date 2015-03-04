package caffe

import (
    "os"
    "os/exec"
)

var CAFFE = "caffe"
var COMPUTER_IMAGE_MEAN = "compute_image_mean"
var CONVERT_IMAGESET = "convert_imageset"

var currentProcess *os.Process

func Run(args... string) error {
    cmd := exec.Command(CAFFE, args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    currentProcess = cmd.Process
    return cmd.Run()
}

func ComputeImageMean(sourceFile, binaryFile string) error {
    cmd := exec.Command(COMPUTER_IMAGE_MEAN, sourceFile, binaryFile)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    currentProcess = cmd.Process
    return cmd.Run()
}

func ConvertImageset(args... string) error {
    cmd := exec.Command(CONVERT_IMAGESET, args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    currentProcess = cmd.Process
    return cmd.Run()
}

func Kill() (error) {
    return currentProcess.Kill()
}
