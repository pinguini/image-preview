package main

import (
	"net/http"

	"github.com/pinguini/image-preview/server"
)

func main() {
	http.HandleFunc("/", server.DefaultHandler)
	http.HandleFunc("/fill/", server.FillHandler)

	http.ListenAndServe(":8080", nil)
}
