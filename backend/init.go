package backend

import (
    "github.com/go-xorm/xorm"
    "log"
    "time"
)

var engine *xorm.Engine

func init() {
    var err error
    if engine, err = xorm.NewEngine(driverName, *sourceName); err != nil {
        log.Fatal(err)
    }
    engine.TZLocation = time.Local
    if err := engine.Sync(File{}, Tag{}, Dataset{}); err != nil {
        log.Fatal(err)
    }
}

func GetEngine() *xorm.Engine {
    return engine
}
