package caffe

import (
    "os"
    "fmt"
    "bufio"
    "regexp"
    "strconv"
    "os/exec"
    "github.com/Lupino/collect/models"
)

var CAFFE = "caffe"
var COMPUTER_IMAGE_MEAN = "compute_image_mean"
var CONVERT_IMAGESET = "convert_imageset"

var iterLine = regexp.MustCompile(`Iteration (?P<iter>\d+)`)
var accuracyLine = regexp.MustCompile(`accuracy = (?P<accuracy>[0-9\.]+)`)
var lrLine = regexp.MustCompile(`lr = (?P<lr>[0-9\.e-]+)`)
var lossLine = regexp.MustCompile(`loss = (?P<loss>[0-9\.]+)`)

func Train(solverFile string) error {
    cmd := exec.Command(CAFFE, "train", "--solver=" + solverFile)

    stderr, _ := cmd.StderrPipe()

    if err := cmd.Start(); err != nil {
        return err
    }

    reader := bufio.NewReader(stderr)
    models.ResetHistory()

    var iter int
    var lr, loss, acc float64

    for {
        line, _, err := reader.ReadLine()
        if err != nil {
            break
        }
        if iterLine.Match(line) {
            iter, _ = strconv.Atoi(string(iterLine.FindSubmatch(line)[1]))
        }
        if lrLine.Match(line) {
            lr, _ = strconv.ParseFloat(string(lrLine.FindSubmatch(line)[1]), 64)
        }
        if lossLine.Match(line) {
            loss, _ = strconv.ParseFloat(string(lossLine.FindSubmatch(line)[1]), 64)
        }
        if accuracyLine.Match(line) {
            acc, _ = strconv.ParseFloat(string(accuracyLine.FindSubmatch(line)[1]), 64)
        }
        models.AddHistory(iter, lr, loss, acc)
        fmt.Printf("%s\n", line)
    }

    cmd.Wait()

    return nil
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
