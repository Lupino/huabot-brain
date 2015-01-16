package main

import (
    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"
    "flag"
)

func init() {
    flag.Parse()
}

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

    api(mart)

    mart.Run()
}
