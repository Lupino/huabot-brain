package models

import (
    "time"
)

type File struct {
    Id        int       `xorm:"pk autoincr"`
    Key       string    `xorm:"varchar(128) notnull unique"`
    CreatedAt time.Time `xorm:"created"`
}

type Tag struct {
    Id        int       `xorm:"pk autoincr"`
    Name      string    `xorm:"varchar(128) notnull unique"`
    CreatedAt time.Time `xorm:"created"`
}


type Dataset struct {
    Id        int       `xorm:"pk autoincr"`
    TagId     int       `xorm:"unique(tag_file)"`
    FileId    int       `xorm:"unique(tag_file)"`
    Type      int       `xorm:"tinyint(1) default(0)"`
    CreatedAt time.Time `xorm:"created"`
}
