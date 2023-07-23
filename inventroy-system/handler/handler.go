package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Golang-Projects/inventory-sytem/model"
	"github.com/Golang-Projects/inventory-sytem/service"
)

type handler struct {
	srv    service.OrderFulfilmentService
	logger *log.Logger
}

type data struct {
	Reserve string `json:"reserve"`
	Message string `json:"message"`
}

type resp struct {
	Code string `json:"code"`
	Data data   `json:"data"`
}

func New(srv service.OrderFulfilmentService, logger *log.Logger) handler {
	return handler{srv: srv, logger: logger}
}

func (h handler) FulFil(w http.ResponseWriter, req *http.Request) {
	reqBytes, _ := ioutil.ReadAll(req.Body)

	var orderReq model.OrderRequest

	err := json.Unmarshal(reqBytes, &orderReq)
	if err != nil {
		h.logger.Println("ERROR", err)
		jsonResp(w, http.StatusBadRequest, resp{Code: "Success", Data: data{"false", "Marshal Error"}})
		return
	}

	if h.srv.CanFulfilOrder(orderReq) {
		jsonResp(w, http.StatusOK, map[string]bool{"can_fulfil": true})

		return
	}

	jsonResp(w, http.StatusOK, map[string]bool{"can_fulfil": false})

}

func (h handler) Reserve(w http.ResponseWriter, req *http.Request) {
	reqBytes, _ := ioutil.ReadAll(req.Body)

	var orderReq model.OrderRequest

	err := json.Unmarshal(reqBytes, &orderReq)
	if err != nil {
		h.logger.Println("ERROR", err)
		jsonResp(w, http.StatusBadRequest, resp{Code: "Success", Data: data{"false", "Marshal Error"}})
		return
	}

	if h.srv.CanFulfilOrder(orderReq) {
		h.srv.ReserveOrder(orderReq)
		jsonResp(w, http.StatusOK, resp{Code: "Success", Data: data{"true", "Success"}})
		return
	}

	jsonResp(w, http.StatusBadRequest, resp{Code: "Success", Data: data{"false", "Insufficient quantities!"}})
}

func jsonResp(w http.ResponseWriter, statusCode int, i interface{}) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(i)
}
