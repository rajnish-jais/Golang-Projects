package store

type Warehouse struct {
	ItemFulfillment map[int]EntityFulfill
}

type EntityFulfill struct {
	ItemName   string
	DateQtyMap map[string]int
}

var wh = map[string]Warehouse{
	"100": {
		map[int]EntityFulfill{
			1: {"Washington Apple (1 pc)", map[string]int{"2021-06-10": 10}},
			2: {"Banana (0.5kg)", map[string]int{"2021-06-10": 10}},
			3: {"Parle-G Biscuit (200g)", map[string]int{"2021-06-10": 10}},
		},
	},
	"200": {
		map[int]EntityFulfill{
			1: {"Washington Apple (1 pc)", map[string]int{"2021-06-10": 10}},
			2: {"Banana (0.5kg)", map[string]int{"2021-06-10": 10}},
			3: {"Parle-G Biscuit (200g)", map[string]int{"2021-06-10": 10}},
		},
	},
}

var cat = map[string]map[string]int{
	"2021-06-10": {"F_N_V": 30, "Grocery": 20},
}
