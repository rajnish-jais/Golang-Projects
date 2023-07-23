package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Golang-Projects/inventory-sytem/handler"
	"github.com/Golang-Projects/inventory-sytem/service"
	"github.com/Golang-Projects/inventory-sytem/store"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	str := store.New()
	svc := service.New(str)
	h := handler.New(svc, logger)

	http.HandleFunc("/fulfil", h.FulFil)
	http.HandleFunc("/reserve", h.Reserve)

	logger.Println("Starting Server")
	logger.Fatal(http.ListenAndServe(":8000", nil))
}
