// +build !mysql

package backend

import (
    _ "github.com/mattn/go-sqlite3"
    "flag"
)

var sourceName = flag.String("dbpath", "resourses/huabot-brain.db", "Sqlite db file.")
var driverName = "sqlite3"
