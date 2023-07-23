package service

import (
	"github.com/Golang-Projects/inventory-sytem/model"
	"github.com/Golang-Projects/inventory-sytem/store"
)

type service struct {
	str store.OrderRequest
}

func New(s store.OrderRequest) OrderFulfilmentService {
	return service{str: s}
}

func (s service) CanFulfilOrder(orderRequest model.OrderRequest) bool {
	sum := make(map[string]int)
	wh, cat := s.str.FetchItemQuantity(orderRequest.WarehouseId, orderRequest.DeliveryDate)

	// item can be fulfil form warehouse or not
	for _, v := range orderRequest.Items {
		if wh.ItemFulfillment[v.ItemId].DateQtyMap[orderRequest.DeliveryDate] < v.Quantity {
			return false
		}
		// sum of required quantity for each category
		sum[v.Category] += v.Quantity
	}

	// item can be fulfil based on category or not
	for k, v := range sum {
		if cat[k] < v {
			return false
		}
	}

	return true
}

func (s service) ReserveOrder(orderRequest model.OrderRequest) {
	s.str.ReserveQuantity(orderRequest)
}
