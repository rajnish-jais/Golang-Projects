package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Golang-Projects/fulfilment-service/handler"
	"github.com/Golang-Projects/fulfilment-service/service"
	"github.com/Golang-Projects/fulfilment-service/store"
	"github.com/gorilla/mux"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	logger.Println("Starting Server")

	store := store.NewInventory()
	service := service.NewInventoryService(store)
	handler := handler.New(service, logger)

	r := mux.NewRouter()

	r.HandleFunc("/addItem", handler.AddItem).Methods("POST")
	r.HandleFunc("/removeItem", handler.RemoveItem).Methods("DELETE")
	r.HandleFunc("/viewItems", handler.ViewItems).Methods("GET")

	logger.Fatal(http.ListenAndServe(":8080", r))
}
