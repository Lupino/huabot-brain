// +build mysql

package backend

import (
    _ "github.com/go-sql-driver/mysql"
    "flag"
)

var sourceName = flag.String("dbpath", "root:@/huabot?charset=utf8", "mysql path.")
var driverName = "mysql"
