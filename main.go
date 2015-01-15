package main

import (
    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"
    _ "github.com/mattn/go-sqlite3"
    "github.com/Lupino/collect/models"
    "github.com/go-xorm/xorm"
    "log"
)

func main() {
    mart := martini.Classic()
    mart.Use(render.Renderer(render.Options{
        Directory: "templates",
        Layout: "layout",
        Extensions: []string{".tmpl", ".html"},
        Charset: "UTF-8",
        IndentJSON: true,
        IndentXML: true,
        HTMLContentType: "application/xhtml+xml",
    }))

    engine, err := xorm.NewEngine("sqlite3", "collect.db")

    if err != nil {
        log.Fatal(err)
    }

    models.Init(engine)

    api(mart, engine)

    mart.Run()
}
