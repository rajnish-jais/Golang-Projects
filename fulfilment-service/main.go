package main

import (
	"github.com/Golang-Projects/fulfilment-service/handler"
	"github.com/Golang-Projects/fulfilment-service/service"
	"github.com/Golang-Projects/fulfilment-service/store"
	"log"
	"net/http"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	logger.Println("Starting Server")

	store := store.NewInventory()
	service := service.NewInventoryService(store)
	handler := handler.New(service, logger)

	http.HandleFunc("/addItem", handler.AddItem)
	http.HandleFunc("/removeItem", handler.RemoveItem)
	http.HandleFunc("/viewItems", handler.ViewItems)

	logger.Fatal(http.ListenAndServe(":8080", nil))
}
