package service

import "Golang-Projects/inventory-sytem/model"

type OrderFulfilmentService interface {
	CanFulfilOrder(orderRequest model.OrderRequest) bool
	ReserveOrder(orderRequest model.OrderRequest)
}
