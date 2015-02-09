package models

import (
    "time"
)

type History struct {
    Id        int       `xorm:"pk autoincr"          json:"file_id,omitempty"`
    Key       string    `xorm:"varchar(128) notnull" json:"key,omitempty"`
    Value     string    `xorm:"varchar(128) notnull" json:"value,omitempty"`
    Timestamp time.Time `xorm:"timestamp"            json:"timestamp,omitempty"`
}

func AddHistory(key, value string) (err error){
    var his = &History{Key: key, Value: value, Timestamp: time.Now()}
    _, err = engine.Insert(his)
    return
}

func FetchHistory() (hist []History, err error) {
    hist = make([]History, 0)
    err = engine.Find(&hist)
    return
}

func ResetHistory() (err error) {
    if err = engine.DropTables(History{}); err != nil {
        return
    }
    return engine.Sync(History{})
}
