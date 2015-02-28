package caffe

import (
    "os"
    "os/exec"
)

var CAFFE = "caffe"
var COMPUTER_IMAGE_MEAN = "compute_image_mean"
var CONVERT_IMAGESET = "convert_imageset"

func Train(solverFile string) error {
    return Run("train", "--solver=" + solverFile)
}

func Run(args... string) error {
    cmd := exec.Command(CAFFE, args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd.Run()
}

func ComputeImageMean(sourceFile, binaryFile string) error {
    cmd := exec.Command(COMPUTER_IMAGE_MEAN, sourceFile, binaryFile)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd.Run()
}

func ConvertImageset(args... string) error {
    cmd := exec.Command(CONVERT_IMAGESET, args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd.Run()
}
