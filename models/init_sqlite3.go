// +build !mysql

package models

import (
    _ "github.com/mattn/go-sqlite3"
    "flag"
)

var sourceName = flag.String("dbpath", "dataset.db", "Sqlite db file.")
var driverName = "sqlite3"
