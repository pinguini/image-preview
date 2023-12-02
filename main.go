package main

import "fmt"

func main() {
    name := "Go Developers"
    fmt.Println("Azure for", name)
}

/*
import (
	"net/http"
	"image-preview/server"
)

func main() {
	http.HandleFunc("/", server.ProxyHandler)
	http.ListenAndServe(":8080", nil)
}
*/