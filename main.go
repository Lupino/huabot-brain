package main

import (
    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"
    "flag"
    "io/ioutil"
    "net/http"
)

func init() {
    flag.Parse()
}

func main() {
    mart := martini.Classic()
    mart.Use(render.Renderer(render.Options{
        Charset: "UTF-8",
        IndentJSON: true,
        IndentXML: true,
        HTMLContentType: "application/xhtml+xml",
    }))

    api(mart)

    mart.Get("/", func(r render.Render) {
        data, _ := ioutil.ReadFile("public/index.html")
        r.Header().Set(render.ContentType, render.ContentHTML)
        r.Data(http.StatusOK, data)
    })

    mart.Run()
}
