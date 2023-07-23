package store

type Inventory interface {
	AddItem(string, int)
	RemoveItem(string) bool
	ViewItem() map[string]int
}
