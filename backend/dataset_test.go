package backend

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

// var engine_test *xorm.Engine

func init() {
    var err error
    engine, err = xorm.NewEngine("sqlite3", "dataset_test.db")
    // engine, err = xorm.NewEngine("mysql", "root:@/collect?charset=utf8")
    if err != nil {
        log.Fatal(err)
    }
    if err := engine.DropTables(File{}, Tag{}, Dataset{}); err != nil {
        log.Fatal(err)
    }
    if err := engine.Sync(File{}, Tag{}, Dataset{}); err != nil {
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

func TestDataset(t *testing.T) {
    var dataset = &Dataset{FileId: 1, TagId: 1}
    if _, err := engine.Insert(dataset); err != nil {
        t.Fatal(err)
    }
    Convey("Dataset: unique(file_tag)", t, func() {
        _, err := engine.Insert(dataset)
        So(err, ShouldNotBeNil)
    })
}

