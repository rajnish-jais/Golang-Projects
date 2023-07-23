package store

var store = map[string]int{
	"product1": 5,
	"product2": 6,
	"product3": 10,
}

type Order struct {
}

func NewStore(order Order) Store {
	return store
}
func addItem(productId string, quantity int) {
	_, exits := store[productId]
	if !exits {
		store[productId] = quantity
	}
	store[productId] += quantity
}

func remove(productId string) bool {
	_, exits := store[productId]
	if !exits {
		return false
	}

	delete(store, productId)
	return false
}

func view() map[string]int {
	return store
}
