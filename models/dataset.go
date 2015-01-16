package models

import (
    "time"
)

const (
    CANDIDATE uint = 0
    TRAIN     uint = 1
    VAL       uint = 2
)

type File struct {
    Id        int       `xorm:"pk autoincr"                 json:"file_id,omitempty"`
    Key       string    `xorm:"varchar(128) notnull unique" json:"key,omitempty"`
    CreatedAt time.Time `xorm:"created"                     json:"created_at,omitempty"`
}

type Tag struct {
    Id        int       `xorm:"pk autoincr"                 json:"tag_id,omitempty"`
    Name      string    `xorm:"varchar(128) notnull unique" json:"name,omitempty"`
    CreatedAt time.Time `xorm:"created"                     json:"created_at,omitempty"`
}


type Dataset struct {
    Id        int       `xorm:"pk autoincr"           json:"dataset_id,omitempty"`
    TagId     int       `xorm:"unique(tag_file)"      json:"tag_id,omitempty"`
    Tag       *Tag      `xorm:"-"                     json:"tag,omitempty"`
    FileId    int       `xorm:"unique(tag_file)"      json:"file_id,omitempty"`
    File      *File     `xorm:"-"                     json:"file,omitempty"`
    DataType  uint      `xorm:"tinyint(1) default(0)" json:"data_type,omitempty"`
    CreatedAt time.Time `xorm:"created"               json:"created_at,omitempty"`
}

func (dataset *Dataset) FillObject() (err error) {
    var tag = &Tag{Id: dataset.TagId}
    var has bool

    if has, err = engine.Get(tag); err != nil {
        return
    } else if has {
        dataset.Tag = tag
    }

    var file = &File{Id: dataset.FileId}
    if has, err = engine.Get(file); err != nil {
        return
    } else if has {
        dataset.File = file
    }

    return
}

func (dataset *Dataset) SetFile(file *File) {
    dataset.File = file
    dataset.FileId = file.Id
}

func (dataset *Dataset) SetTag(tag *Tag) {
    dataset.Tag = tag
    dataset.TagId = tag.Id
}
