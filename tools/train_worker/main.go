package main

import (
    "github.com/Lupino/huabot-brain/tools/train_worker/caffe"
    "github.com/mikespook/gearman-go/worker"
    "net/http"
    "log"
    "io"
    "os"
    "os/exec"
    "io/ioutil"
    "sync"
    "flag"
)

var trainLock = false
var trainLocker = new(sync.Mutex)


var (
    resoursesPath = flag.String("resourses", "resourses", "The resourses path.")
    UPLOADPATH = flag.String("datasets", "public/upload/", "The datasets path.")
    TRAIN_FILE = *resoursesPath + "/train.txt"
    VAL_FILE = *resoursesPath + "/val.txt"
    TRAIN_LMDB = *resoursesPath + "/train_lmdb"
    VAL_LMDB = *resoursesPath + "/val_lmdb"
    MEAN_FILE = *resoursesPath + "/mean.binaryproto"
    SOLVER_FILE = *resoursesPath + "/solver.prototxt"
    LOG_DIR = *resoursesPath + "/logs"
    PLOT_ROOT = *resoursesPath + "/plot"
    BRAIN_ROOT = flag.String("api", "http://127.0.0.1:3000/api", "The Huabot Brain api root.")
    GEARMAND = flag.String("gearmand", "127.0.0.1:4730", "The Gearmand server.")
)

func init() {
    flag.Parse()
}

func loadFile(url string, file *os.File) (err error) {
    log.Printf("load %s\n", url)
    var resp *http.Response
    if resp, err = http.Get(url); err != nil {
        return
    }
    defer resp.Body.Close()
    io.Copy(file, resp.Body)
    return
}

func loadDataset(host string) (err error) {

    var trainFile, valFile *os.File
    if trainFile, err = os.Create(TRAIN_FILE); err != nil {
        return
    }
    defer trainFile.Close()

    if err = loadFile(host + "/train.txt", trainFile); err != nil {
        return
    }

    if valFile, err = os.Create(VAL_FILE); err != nil {
        return
    }
    defer valFile.Close()

    if err = loadFile(host + "/val.txt", valFile); err != nil {
        return
    }

    return
}

func caffeTrain() {
    var err error
    log.Printf("caffeTrain")
    if trainLock {
        return
    }

    trainLock = true
    trainLocker.Lock()
    defer (func() {
        trainLock = false
        trainLocker.Unlock()
    })()

    if err = loadDataset(*BRAIN_ROOT); err != nil {
        log.Printf("Error: %s\n", err)
        return
    }

    os.RemoveAll(TRAIN_LMDB)
    os.RemoveAll(VAL_LMDB)


    if err = caffe.ConvertImageset("--resize_height=256", "--shuffle", "--resize_width=256", *UPLOADPATH, TRAIN_FILE, TRAIN_LMDB); err != nil {
        log.Printf("Error: %s\n", err)
        return
    }

    if err = caffe.ConvertImageset("--resize_height=256", "--resize_width=256", "--shuffle", *UPLOADPATH, VAL_FILE, VAL_LMDB); err != nil {
        log.Printf("Error: %s\n", err)
        return
    }

    if err = caffe.ComputeImageMean(TRAIN_LMDB, MEAN_FILE); err != nil {
        log.Printf("Error: %s\n", err)
        return
    }

    if err = caffe.Run("train", "--solver=" + SOLVER_FILE, "-log_dir=" + LOG_DIR); err != nil {
        log.Printf("Error: %s\n", err)
        return
    }
}

func caffeTrainTask(job worker.Job) (data []byte, err error) {
    go caffeTrain()
    return []byte("training"), nil
}

func caffeStatusTask(job worker.Job) (data []byte, err error) {
    status := "no training"
    if trainLock {
        status = "training"
    }
    return []byte(status), nil
}

func caffePlotTask(job worker.Job) (data []byte, err error) {
    suffix := string(job.Data())

    if err = run(PLOT_ROOT + "/parse_log.sh", LOG_DIR + "/caffe.INFO"); err != nil {
        log.Printf("Error: %s\n", err)
        return
    }

    if err = run("gnuplot", PLOT_ROOT + "/plot_log.gnuplot." + suffix); err != nil {
        log.Printf("Error: %s\n", err)
        return
    }
    data, err = ioutil.ReadFile("/tmp/" + suffix + ".png")
    return
}

func run(cmdName string, args... string) (error) {
    cmd := exec.Command(cmdName, args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd.Run()
}

func main() {
    w := worker.New(worker.OneByOne)
    w.ErrorHandler = func(e error) {
        log.Println(e)
    }
    w.AddServer("tcp4", *GEARMAND)
    w.AddFunc("CAFFE:TRAIN", caffeTrainTask, worker.Unlimited)
    w.AddFunc("CAFFE:TRAIN:STATUS", caffeStatusTask, worker.Unlimited)
    w.AddFunc("CAFFE:TRAIN:PLOT", caffePlotTask, worker.Unlimited)

    if err := w.Ready(); err != nil {
        log.Fatal(err)
        return
    }
    w.Work()
}
