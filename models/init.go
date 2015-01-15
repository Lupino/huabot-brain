package models

import (
    "github.com/go-xorm/xorm"
    "log"
)

func Init(engine *xorm.Engine) {
    if err := engine.Sync(File{}, Tag{}, Dataset{}); err != nil {
        log.Fatal(err)
    }
}
