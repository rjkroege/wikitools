package main

import (
    "fmt"
    "log"
    "net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "<h1>Hello, World!</h1>")
}

func main() {
    http.HandleFunc("/", hello)
    log.Println("Listening on http://localhost:8080 …")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
