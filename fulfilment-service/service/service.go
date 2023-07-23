package service

import "github.com/Golang-Projects/fulfilment-service/store"

type service struct {
	inv store.Inventory
}

func NewInventoryService(inventory store.Inventory) InventoryService {
	return service{inv: inventory}
}

func (srv service) AddItemToInventory(productID string, quantity int) {
	srv.inv.AddItem(productID, quantity)
}

func (srv service) RemoveItemFromInventory(productID string) bool {
	return srv.inv.RemoveItem(productID)
}

func (srv service) ViewItemFromInventory() map[string]int {
	return srv.inv.ViewItem()
}
