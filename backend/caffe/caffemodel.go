package caffe

import (
    "os"
    "path/filepath"
    "github.com/Lupino/huabot-brain/config"
)

func ListModels() (modelNames []string, err error) {
    if modelNames, err = filepath.Glob(config.CAFFEMODEL_PATH + "/*.caffemodel"); err != nil {
        return
    }
    for i, modelName := range modelNames {
        modelNames[i] = filepath.Base(modelName)
    }
    return
}

func ApplyModel(modelName string) (error) {
    os.Remove(config.CAFFEMODEL_NAME)
    return os.Symlink(filepath.Join(config.CAFFEMODEL_PATH, modelName), config.CAFFEMODEL_NAME)
}

func RemoveModel(modelName string) (error) {
    return os.Remove(filepath.Join(config.CAFFEMODEL_PATH, modelName))
}

func GetCurrentModel() (modelName string, err error) {
    if modelName, err = os.Readlink(config.CAFFEMODEL_NAME); err != nil {
        return
    }
    modelName = filepath.Base(modelName)
    return
}
