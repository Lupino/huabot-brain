// +build mysql

package models

import (
    _ "github.com/go-sql-driver/mysql"
    "flag"
)

var sourceName = flag.String("dbpath", "root:@/collect?charset=utf8", "mysql path.")
var driverName = "mysql"
