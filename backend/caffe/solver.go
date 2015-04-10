package caffe

import (
    "os"
    "fmt"
    "log"
    "sync"
    "github.com/Lupino/huabot-brain/backend"
)

var (
    resoursesPath = "resourses"
    UPLOADPATH = "public/upload/"
    TRAIN_FILE = resoursesPath + "/train.txt"
    VAL_FILE = resoursesPath + "/val.txt"
    TRAIN_LMDB = resoursesPath + "/train_lmdb"
    VAL_LMDB = resoursesPath + "/val_lmdb"
    MEAN_FILE = resoursesPath + "/mean.binaryproto"
    SOLVER_FILE = resoursesPath + "/solver.prototxt"
    LOG_DIR = resoursesPath + "/logs"
)

var solverLocker = new(sync.Mutex)
var solveLocked = false

func exportToFile(file *os.File, dataType uint) (err error) {
    var engine = backend.GetEngine()
    err = engine.Where("data_type=?", dataType).Iterate(new(backend.Dataset),
                                                        func(i int, bean interface{}) error {
        dataset := bean.(*backend.Dataset)
        dataset.FillObject()
        var ext, ok = backend.FILE_EXTS[dataset.File.Type]
        if !ok {
            ext = ".jpg"
        }
        fmt.Fprintf(file, "%s%s %d\n", dataset.File.Key, ext, dataset.TagId)
        return nil
    })
    return
}

func prepare() (err error) {
    var trainFile, valFile *os.File
    if trainFile, err = os.Create(TRAIN_FILE); err != nil {
        return
    }
    defer trainFile.Close()

    if !IsOnSolving() {
        return
    }

    if err = exportToFile(trainFile, backend.TRAIN); err != nil {
        return
    }

    if !IsOnSolving() {
        return
    }

    if valFile, err = os.Create(VAL_FILE); err != nil {
        return
    }
    defer valFile.Close()

    if !IsOnSolving() {
        return
    }

    if err = exportToFile(valFile, backend.VAL); err != nil {
        return
    }

    if !IsOnSolving() {
        return
    }

    os.RemoveAll(TRAIN_LMDB)
    os.RemoveAll(VAL_LMDB)

    if !IsOnSolving() {
        return
    }

    if err = ConvertImageset("--resize_height=256",
                             "--resize_width=256",
                             "--shuffle",
                             UPLOADPATH,
                             TRAIN_FILE, TRAIN_LMDB); err != nil {
        return
    }

    if err = ConvertImageset("--resize_height=256",
                             "--resize_width=256",
                             "--shuffle",
                             UPLOADPATH,
                             VAL_FILE, VAL_LMDB); err != nil {
        return
    }

    if err = ComputeImageMean(TRAIN_LMDB, MEAN_FILE); err != nil {
        return
    }
    return
}

func Solver() {
    var err error

    if solveLocked {
        return
    }

    solverLocker.Lock()
    solveLocked = true
    defer (func() {
        solveLocked = false
        solverLocker.Unlock()
    })()

    if err = prepare(); err != nil {
        log.Printf("Error: prepare dataset fail.")
        return
    }

    if !IsOnSolving() {
        return
    }

    if err = Run("train", "--solver=" + SOLVER_FILE, "-log_dir=" + LOG_DIR); err != nil {
        log.Printf("Error: %s\n", err)
        return
    }
    return
}

func Solve() {
    go Solver()
}

func IsOnSolving() bool {
    return solveLocked
}

func StopSolve() {
    Kill()
    solveLocked = false
}
