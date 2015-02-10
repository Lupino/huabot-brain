package caffe

import (
    "os"
    "fmt"
    "bufio"
    "regexp"
    "os/exec"
    "github.com/Lupino/collect/models"
)

var CAFFE = "caffe"
var COMPUTER_IMAGE_MEAN = "compute_image_mean"
var CONVERT_IMAGESET = "convert_imageset"

var RRD_FILE = "dataset.rrd"
var STEP uint = 1
var HEARTBEAT = 2 * STEP

var iterLine = regexp.MustCompile(`Iteration (?P<iter>\d+)`)
var accuracyLine = regexp.MustCompile(`accuracy = (?P<accuracy>[0-9\.]+)`)
var lrLine = regexp.MustCompile(`lr = (?P<lr>[0-9\.e-]+)`)
var lossLine = regexp.MustCompile(`loss = (?P<loss>[0-9\.]+)`)

func Train(solverFile string) error {
    return Run("train", "--solver=" + solverFile)
}

func Run(args... string) error {
    cmd := exec.Command(CAFFE, args...)

    stderr, _ := cmd.StderrPipe()

    if err := cmd.Start(); err != nil {
        return err
    }

    reader := bufio.NewReader(stderr)
    models.ResetHistory()

    for {
        line, _, err := reader.ReadLine()
        if err != nil {
            break
        }
        if iterLine.Match(line) {
            models.AddHistory("iter", string(iterLine.FindSubmatch(line)[1]))
        }
        if lrLine.Match(line) {
            models.AddHistory("lr", string(lrLine.FindSubmatch(line)[1]))
        }
        if lossLine.Match(line) {
            models.AddHistory("loss", string(lossLine.FindSubmatch(line)[1]))
        }
        if accuracyLine.Match(line) {
            models.AddHistory("acc", string(accuracyLine.FindSubmatch(line)[1]))
        }
        fmt.Printf("%s\n", line)
    }

    cmd.Wait()

    return nil
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
