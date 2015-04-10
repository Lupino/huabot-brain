package caffe

import (
    "os"
    "fmt"
    "log"
    "sync"
    "github.com/Lupino/huabot-brain/backend"
    "github.com/Lupino/huabot-brain/config"
)

var solverLocker = new(sync.Mutex)
var solveLocked = false

func exportToFile(file *os.File, dataType uint) (err error) {
    var engine = backend.GetEngine()
    err = engine.Where("data_type=?", dataType).Iterate(new(backend.Dataset),
                                                        func(i int, bean interface{}) error {
        dataset := bean.(*backend.Dataset)
        dataset.FillObject()
        var ext, ok = config.FILE_EXTS[dataset.File.Type]
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
    if trainFile, err = os.Create(config.TRAIN_FILE); err != nil {
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

    if valFile, err = os.Create(config.VAL_FILE); err != nil {
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

    os.RemoveAll(config.TRAIN_LMDB)
    os.RemoveAll(config.VAL_LMDB)

    if !IsOnSolving() {
        return
    }

    if err = ConvertImageset("--resize_height=256",
                             "--resize_width=256",
                             "--shuffle",
                             config.UPLOADPATH,
                             config.TRAIN_FILE, config.TRAIN_LMDB); err != nil {
        return
    }

    if err = ConvertImageset("--resize_height=256",
                             "--resize_width=256",
                             "--shuffle",
                             config.UPLOADPATH,
                             config.VAL_FILE, config.VAL_LMDB); err != nil {
        return
    }

    if err = ComputeImageMean(config.TRAIN_LMDB, config.MEAN_FILE); err != nil {
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

    if err = Run("train", "--solver=" + config.SOLVER_FILE,
                 "-log_dir=" + config.LOG_DIR); err != nil {
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
