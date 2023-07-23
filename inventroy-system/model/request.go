package model


type ItemRequest struct {
	ItemId   int `json:"item_id"`
	ItemName string `json:"item_name"`
	Category string `json:"category"`
	Quantity int `json:"quantity"`
}

type OrderRequest struct {
	CustomerId   string `json:"customer_id"`
	WarehouseId  string `json:"warehouse_id"`
	DeliveryDate string `json:"delivery_date"`
	Items        []ItemRequest `json:"items"`
}

