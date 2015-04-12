package main

import (
    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"
    "github.com/Lupino/huabot-brain/backend"
    "github.com/Lupino/huabot-brain/config"
    "flag"
    "io/ioutil"
    "net/http"
)

var resourceName = flag.String("resources", "resourses", "Resources path.")
var predictRoot = flag.String("predict", "http://127.0.0.1:3001", "Predict host.")

func init() {
    flag.Parse()
}

func main() {
    config.SetResource(*resourceName)
    config.SetPredictRoot(*predictRoot)

    mart := martini.Classic()
    mart.Use(render.Renderer(render.Options{
        Charset: "UTF-8",
        IndentJSON: true,
        IndentXML: true,
        HTMLContentType: "application/xhtml+xml",
    }))

    backend.Init()

    api(mart)

    mart.Get("/", func(r render.Render) {
        data, _ := ioutil.ReadFile("public/index.html")
        r.Header().Set(render.ContentType, render.ContentHTML)
        r.Data(http.StatusOK, data)
    })

    mart.Run()
}
