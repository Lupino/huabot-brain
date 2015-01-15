package models

import (
    // _ "github.com/go-sql-driver/mysql"
    _ "github.com/mattn/go-sqlite3"
    "github.com/go-xorm/xorm"
    . "github.com/smartystreets/goconvey/convey"
    "testing"
    "strconv"
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
        var file = &File{Key: key}
        if _, err := engine.Insert(file); err != nil {
            t.Fatal(err)
        }
        Convey("File: " + key + "'s Id should equal " + strconv.Itoa(i + 1), t, func() {
            So(file.Id, ShouldEqual, i + 1)
        })
    }
}

func TestFileTag(t *testing.T) {
    var fileTag = &FileTag{FileId: 1, TagId: 1}
    if _, err := engine.Insert(fileTag); err != nil {
        t.Fatal(err)
    }
    Convey("FileTag: unique(file_tag)", t, func() {
        _, err := engine.Insert(fileTag)
        So(err, ShouldNotBeNil)
    })
}

