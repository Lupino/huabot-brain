package models

import (
    // _ "github.com/go-sql-driver/mysql"
    _ "github.com/mattn/go-sqlite3"
    "github.com/go-xorm/xorm"
    "testing"
    "log"
    "fmt"
)

import (
)

var engine *xorm.Engine

func init() {
    var err error
    engine, err = xorm.NewEngine("sqlite3", "file_test.db")
    // engine, err = xorm.NewEngine("mysql", "root:@/collect?charset=utf8")
    if err != nil {
        log.Fatal(err)
    }
    if err := engine.DropTables(File{}, Tag{}, FileTag{}); err != nil {
        log.Fatal(err)
    }
    if err := engine.Sync(File{}, Tag{}, FileTag{}); err != nil {
        log.Fatal(err)
    }
}

func TestInsertFile(t *testing.T) {
    for i := range make([]int, 5) {
        var key = fmt.Sprintf("key %d", i)
        if _, err := engine.Insert(&File{Key: key}); err != nil {
            t.Fatal(err)
        }
    }
}

func TestFileTag(t *testing.T) {
    if _, err := engine.Insert(&FileTag{FileId: 1, TagId: 1}); err != nil {
        t.Fatal(err)
    }
    if _, err := engine.Insert(&FileTag{FileId: 1, TagId: 1}); err == nil {
        t.Fatal("unique key: file_tag fail")
    }
}
