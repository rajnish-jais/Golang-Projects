package handler

import (
	"encoding/json"
	"fmt"
	"github.com/Golang-Projects/fulfilment-service/service"
	"log"
	"net/http"
	"strconv"
)

type handler struct {
	service service.InventoryService
	logger  *log.Logger
}

type resp struct {
	Code  string      `json:"code"`
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

func New(inventoryService service.InventoryService, logger *log.Logger) handler {
	return handler{inventoryService, logger}
}

func (h handler) AddItem(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	productID := query.Get("productID")
	amount := query.Get("amount")

	if productID == "" || amount == "" {
		jsonResp(w, http.StatusBadRequest, resp{Code: "Failure", Error: "Missing Params"})
		return
	}

	quantity, err := strconv.Atoi(amount)
	if err != nil {
		jsonResp(w, http.StatusBadRequest, resp{Code: "Failure", Error: "Incorrect Amount Param"})
		return
	}

	h.service.AddItemToInventory(productID, quantity)

	jsonResp(w, http.StatusOK, resp{Code: "Success"})
}

func (h handler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	productID := query.Get("productID")
	if productID == "" {
		jsonResp(w, http.StatusBadRequest, resp{Code: "Failure", Error: "Missing Params"})
		return
	}

	removed := h.service.RemoveItemFromInventory(productID)
	if !removed {
		jsonResp(w, http.StatusBadRequest, resp{Code: "Success", Error: fmt.Sprintf("%v is not found in the inventory", productID)})
		return
	}
	jsonResp(w, http.StatusOK, resp{Code: "Success", Data: fmt.Sprintf("One item form %v has been removed from inventory", productID)})
}

func (h handler) ViewItems(w http.ResponseWriter, r *http.Request) {
	jsonResp(w, http.StatusOK, resp{Code: "Success", Data: h.service.ViewItemFromInventory()})
}

func jsonResp(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
