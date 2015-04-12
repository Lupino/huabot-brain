package config

var (
    RES string
    UPLOADPATH = "public/upload/"
    TRAIN_FILE string
    VAL_FILE string
    TRAIN_LMDB string
    VAL_LMDB string
    MEAN_FILE string
    SOLVER_FILE string
    LOG_DIR string
    PLOT_ROOT string
    PREDICT_HOST string
    PREDICT string
)

var FILE_EXTS = map[string]string{
    "image/png": ".png",
    "image/jpeg": ".jpg",
    "image/gif": ".gif",
}

func SetResource(source string) {
    RES = source
    TRAIN_FILE = RES + "/train.txt"
    VAL_FILE = RES + "/val.txt"
    TRAIN_LMDB = RES + "/train_lmdb"
    VAL_LMDB = RES + "/val_lmdb"
    MEAN_FILE = RES + "/mean.binaryproto"
    SOLVER_FILE = RES + "/solver.prototxt"
    LOG_DIR = RES + "/logs"
    PLOT_ROOT = RES + "/plot"
    PREDICT = RES + "/predict/main.py"
}

func SetPredictRoot(root string) {
    PREDICT_HOST = root
}
