package main

import (
	"log"
	"net/http"
)

const (
	port = ":3456"
)

func main() {
	dal := NewDAL("postgres://postgres:secret@localhost:5432/steakhouse?sslmode=disable")
	router := CreateRouter(dal)
	log.Fatal(http.ListenAndServe(port, router))
}
