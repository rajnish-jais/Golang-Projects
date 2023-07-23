package service

type InventoryService interface {
	AddItemToInventory(string, int)
	RemoveItemFromInventory(string) bool
	ViewItemFromInventory() map[string]int
}
