package store

var store = map[string]int{
	"product1": 5,
	"product2": 6,
	"product3": 10,
}

type ItemReq struct {
	store map[string]int
}

func NewInventory() Inventory {
	return ItemReq{store: store}
}

func (ir ItemReq) AddItem(productId string, quantity int) {
	ir.store[productId] += quantity
}

func (ir ItemReq) RemoveItem(productId string) bool {
	count, exits := ir.store[productId]
	if !exits {
		return false
	}

	if count == 1 {
		delete(store, productId)
	} else {
		ir.store[productId] = count - 1
	}

	return true
}

func (ir ItemReq) ViewItem() map[string]int {
	return ir.store
}
