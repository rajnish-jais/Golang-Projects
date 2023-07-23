package store

import "github.com/Golang-Projects/inventory-sytem/model"

type order struct {
	wh  map[string]Warehouse
	cat map[string]map[string]int
}

func New() OrderRequest {
	return order{wh: wh, cat: cat}
}

func (o order) FetchItemQuantity(warehouseId, deliveryDate string) (Warehouse, map[string]int) {
	return o.wh[warehouseId], o.cat[deliveryDate]
}

func (o order) ReserveQuantity(orderRequest model.OrderRequest) {
	// updating the items quantity in warehouse and category wise
	for _, v := range orderRequest.Items {
		o.wh[orderRequest.WarehouseId].ItemFulfillment[v.ItemId].DateQtyMap[orderRequest.DeliveryDate] -= v.Quantity
		o.cat[orderRequest.DeliveryDate][v.Category] -= v.Quantity
	}
}
