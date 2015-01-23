package main

import (
    . "github.com/smartystreets/goconvey/convey"
    "github.com/Lupino/collect/models"
    "testing"
    "net/http"
    "os"
    "bytes"
    "path/filepath"
    "mime/multipart"
    "io"
    "encoding/json"
)

var host = "http://127.0.0.1:3000"

// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)
    part, err := writer.CreateFormFile(paramName, filepath.Base(path))
    if err != nil {
        return nil, err
    }
    _, err = io.Copy(part, file)

    for key, val := range params {
        _ = writer.WriteField(key, val)
    }
    err = writer.Close()
    if err != nil {
        return nil, err
    }
    req, err := http.NewRequest("POST", uri, body)
    req.Header.Add("Content-Type", writer.FormDataContentType())
    return req, err
}

func TestPostDataset(t *testing.T) {
    Convey("POST /api/datasets/", t, func() {
        path, _ := os.Getwd()
        path += "/public/favicon.ico"
        extraParams := map[string]string{
            "tag": "test 1",
        }

        request, err := newfileUploadRequest(host + API + "/datasets", extraParams, "file", path)
        So(err, ShouldBeNil)

        client := new(http.Client)
        resp, err := client.Do(request)
        So(err, ShouldBeNil)

        body := &bytes.Buffer{}
        _, err = body.ReadFrom(resp.Body)
        So(err, ShouldBeNil)
        resp.Body.Close()
        var data = make(map[string]models.Dataset)
        json.Unmarshal(body.Bytes(), &data)
        var dataset = data["dataset"]
        So(dataset.Tag.Name, ShouldEqual, "test 1")
    })
}
