package store

import "github.com/Golang-Projects/inventory-sytem/model"

type OrderRequest interface {
	FetchItemQuantity(string, string) (Warehouse, map[string]int)
	ReserveQuantity(model.OrderRequest)
}
