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

var cacheHist = new(History)

func AddHistory(iter int, lr, loss, acc float64) (err error){
    var hist = &History{
        Iter: iter,
        Lr: lr,
        Loss: loss,
        Acc: acc,
        Timestamp: time.Now(),
    }

    if cacheHist.Timestamp.Add(1 * time.Second).Unix() < hist.Timestamp.Unix() {
        _, err = engine.Insert(cacheHist)
    }
    cacheHist = hist
    return
}

func FetchHistory() (hist []History, err error) {
    hist = make([]History, 0)
    err = engine.Find(&hist)
    return
}

func ResetHistory() (err error) {
    cacheHist = new(History)
    cacheHist.Timestamp = time.Now()
    if err = engine.DropTables(History{}); err != nil {
        return
    }
    return engine.Sync(History{})
}
