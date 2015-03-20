package backend

import (
    "io"
    "os"
    "fmt"
    "time"
    "image"
    "crypto/sha1"
    _ "image/png"
    _ "image/jpeg"
    "encoding/hex"
    "mime/multipart"
)

const (
    CANDIDATE uint = 0
    TRAIN     uint = 1
    VAL       uint = 2
    UPLOADPATH = "public/upload/"
)

type File struct {
    Id        int       `xorm:"pk autoincr"                 json:"file_id,omitempty"`
    Key       string    `xorm:"varchar(128) notnull unique" json:"key,omitempty"`
    Width     int       `xorm:"integer(4)"                  json:"width,omitempty"`
    Height    int       `xorm:"integer(4)"                  json:"height,omitempty"`
    CreatedAt time.Time `xorm:"created"                     json:"created_at,omitempty"`
}

type Tag struct {
    Id         int       `xorm:"pk autoincr"                 json:"tag_id,omitempty"`
    Name       string    `xorm:"varchar(128) notnull unique" json:"name,omitempty"`
    TrainCount int       `xorm:"tinyint(1) default(0)"       json:"train_count,omitempty"`
    TestCount  int       `xorm:"tinyint(1) default(0)"       json:"test_count,omitempty"`
    CreatedAt  time.Time `xorm:"created"                     json:"created_at,omitempty"`
}


type Dataset struct {
    Id          int       `xorm:"pk autoincr"           json:"dataset_id,omitempty"`
    TagId       int       `xorm:"unique(tag_file)"      json:"tag_id,omitempty"`
    Tag         *Tag      `xorm:"-"                     json:"tag,omitempty"`
    FileId      int       `xorm:"unique(tag_file)"      json:"file_id,omitempty"`
    File        *File     `xorm:"-"                     json:"file,omitempty"`
    DataType    uint      `xorm:"tinyint(1) default(0)" json:"data_type,omitempty"`
    Description string    `xorm:"varchar(256)"          json:"description,omitempty"`
    CreatedAt   time.Time `xorm:"created"               json:"created_at,omitempty"`
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

func UploadFile(realFile *multipart.FileHeader) (file *File, err error) {
    var source multipart.File
    if source, err = realFile.Open(); err != nil {
        return
    }
    defer source.Close()
    var img image.Config
    if img, _, err = image.DecodeConfig(source); err != nil {
        return
    }
    source.Seek(0, 0)
    hasher := sha1.New()
    io.Copy(hasher, source)
    fileKey := hex.EncodeToString(hasher.Sum(nil))
    file = &File{Key: fileKey}
    has, _ := engine.Get(file)
    if !has {
        var dst *os.File
        if dst, err = os.Create(UPLOADPATH + fileKey); err != nil {
            return
        }
        defer dst.Close()
        source.Seek(0, 0)
        if _, err = io.Copy(dst, source); err != nil {
            return
        }

        file.Width = img.Width
        file.Height = img.Height

        if _, err = engine.Insert(file); err != nil {
            return
        }
    }
    return
}

func SaveTag(realTag string) (tag *Tag, err error) {
    tag = &Tag{Name: realTag}
    has, _ := engine.Get(tag)
    if !has {
        if _, err = engine.Insert(tag); err != nil {
            return
        }
    }
    return
}

func DeleteTag(tagId int) (err error) {
    var q = engine.Where("tag_id = ?", tagId)
    var dataset Dataset
    var tag Tag
    q.Delete(&dataset)
    engine.Id(tagId).Delete(&tag)
    return
}

func SaveDataset(file *File, tag *Tag, dataType uint, desc string) (dataset *Dataset, err error) {
    dataset = &Dataset{FileId: file.Id, TagId: tag.Id}
    has, _ := engine.Get(dataset)
    if !has {
        dataset.DataType = dataType
        dataset.Description = desc
        if _, err = engine.Insert(dataset); err != nil {
            return
        }
        var sql string
        if dataType == TRAIN {
          sql = "update `tag` set `train_count` = `train_count` + 1 where `id` = ?"
        } else if dataType == VAL {
          sql = "update `tag` set `test_count` = `test_count` + 1 where `id` = ?"
        }
        engine.Exec(sql, tag.Id)
    }
    dataset.File = file
    dataset.Tag = tag
    return
}

func ExportDataset(dataType uint) (text string, err error) {
    err = engine.Where("data_type=?", dataType).Iterate(new(Dataset), func(i int, bean interface{}) error {
        dataset := bean.(*Dataset)
        dataset.FillObject()
        text = fmt.Sprintf("%s%s %d\n", text, dataset.File.Key, dataset.TagId)
        return nil
    })
    return
}

