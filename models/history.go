package models

import (
    "time"
)

type History struct {
    Id        int       `xorm:"pk autoincr"     json:"hist_id,omitempty"`
    Iter      int       `xorm:"int(10) notnull" json:"iter,omitempty"`
    Lr        float64   `                       json:"lr,omitempty"`
    Loss      float64   `                       json:"loss,omitempty"`
    Acc       float64   `                       json:"acc,omitempty"`
    Timestamp time.Time `xorm:"timestamp"       json:"timestamp,omitempty"`
}

func AddHistory(iter int, lr, loss, acc float64) (err error){
    var his = &History{
        Lr: lr,
        Loss: loss,
        Acc: acc,
        Timestamp: time.Now(),
    }
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
