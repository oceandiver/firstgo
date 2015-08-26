package main

import (
  "github.com/codegangsta/negroni"
  "github.com/gorilla/mux"
  session "github.com/oceandiver/letsgo/lib/session"
  "net/http"
  "path"
  "fmt"
  "os"
)

func main() {
  fmt.Println("Server started!")

  r := mux.NewRouter()

  // static files
  root, _ := os.Getwd()
  r.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir(path.Join(root, "public")))))

// router goes last
  /*r.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(w, "Welcome to the home page!")
  }) */

  r.PathPrefix("/").Handler(session.BuildRoutes())

  n := negroni.New()
  n.UseHandler(r)
  n.Run(":8080")

}
