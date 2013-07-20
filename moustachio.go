package main

import (
  "code.google.com/p/freetype-go/freetype/raster"
  "html/template"
  "io"
  "io/ioutil"
  "net/http"
)

var uploadTemplate = template.Must(template.ParseFiles("upload.html"))
var errorTemplate = template.Must(template.ParseFiles("error.html"))

func errorHandler(fn http.HandlerFunc) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    defer func() {
      if e := recover(); e != nil {
        w.WriteHeader(500)
        errorTemplate.Execute(w, e)
      }
    }()
    fn(w, r)
  }
}

func check(err error) {
  if err != nil {
    panic(err)
  }
}

func upload(w http.ResponseWriter, r *http.Request) {
  if r.Method != "POST" {
    uploadTemplate.Execute(w, nil)
    return
  }

  // opens sent file
  f, _, err := r.FormFile("image")
  check(err)
  defer f.Close()

  // creates a new temp file
  t, err := ioutil.TempFile(".", "image-")
  check(err)
  defer t.Close()

  // copy the content of the file to disk
  _, err = io.Copy(t, f)
  check(err)
  http.Redirect(w, r, "/view?id="+t.Name()[6:], 302)
}

func view(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "image")
  http.ServeFile(w, r, "image-"+r.FormValue("id"))
}

func main() {
  http.HandleFunc("/", errorHandler(upload))
  http.HandleFunc("/view", errorHandler(view))
  http.ListenAndServe(":8080", nil)
}
