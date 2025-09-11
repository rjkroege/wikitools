package main

import (
    "flag"
    "fmt"
    "log"
    "net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "<h1>Hello, World!</h1>")
}

func main() {
    port := flag.String("port", "8080", "HTTP listen port")
    flag.Parse()

    addr := ":" + *port
    http.HandleFunc("/", hello)
    log.Printf("Listening on http://localhost%s …", addr)
    log.Fatal(http.ListenAndServe(addr, nil))
}
