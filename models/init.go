package models

import (
    _ "github.com/mattn/go-sqlite3"
    "github.com/go-xorm/xorm"
    "flag"
    "log"
)

var sourceName = flag.String("dbpath", "dataset.db", "Sqlite db file.")

var engine *xorm.Engine

func init() {
    var err error
    if engine, err = xorm.NewEngine("sqlite3", *sourceName); err != nil {
        log.Fatal(err)
    }
    if err := engine.Sync(File{}, Tag{}, Dataset{}); err != nil {
        log.Fatal(err)
    }
}

func GetEngine() *xorm.Engine {
    return engine
}
