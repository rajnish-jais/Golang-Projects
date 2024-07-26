// https://medium.com/@interviewready/low-level-design-of-splitwise-f334c8f6ff77
/*
Time Complexity
Adding, Editing, and Settling Expenses:O(1) for each operation since these operations involve basic insert, update, or lookup operations in a map.

Creating Groups: O(1) for creating a group since it involves inserting a new entry into the map.

Generating Payment Graph: O(nlogn), where n is the number of users with non-zero balances. This involves heap operations for dividing users into positive and negative balances and adjusting the balances.

Space Complexity
Users, Groups, and Expenses Storage:O(u+g+e), where
u is the number of users, g is the number of groups, and e is the number of expenses.
Balancing Algorithm:O(n) for storing user balances in heaps, where n is the number of users with non-zero balances.
This implementation provides the core functionalities for managing expenses, editing them, settling them, creating groups, and generating a payment graph to minimize the number of transactions.
*/

package main

import (
	"container/heap"
	"fmt"
	"time"
)

// User represents a user with an ID, image URI, and bio.
type User struct {
	UserID   string
	ImageURI string
	Bio      string
}

// Balance represents the balance of a user in a specific currency.
type Balance struct {
	Currency string
	Amount   int
}

// Expense represents an expense with various metadata and balance information.
type Expense struct {
	ExpenseID string
	IsSettled bool
	Balances  map[string]Balance
	GroupID   string
	Title     string
	Timestamp int64
	ImageURI  string
}

// Group represents a group of users with various metadata.
type Group struct {
	GroupID     string
	Users       []User
	ImageURI    string
	Title       string
	Description string
}

// PaymentNode represents a payment between two users.
type PaymentNode struct {
	From   string
	To     string
	Amount int
}

// MaxHeap is a max-heap of Nodes.
type MaxHeap []Node

func (h MaxHeap) Len() int           { return len(h) }
func (h MaxHeap) Less(i, j int) bool { return h[i].FinalBalance > h[j].FinalBalance }
func (h MaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MaxHeap) Push(x interface{}) {
	*h = append(*h, x.(Node))
}

func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// Node represents a user node with their final balance.
type Node struct {
	UserID       string
	FinalBalance int
}

var (
	users    = make(map[string]User)
	groups   = make(map[string]Group)
	expenses = make(map[string]Expense)
)

func makePaymentGraph(balances map[string]int) []PaymentNode {
	var posHeap, negHeap MaxHeap

	for uid, balance := range balances {
		if balance > 0 {
			heap.Push(&posHeap, Node{UserID: uid, FinalBalance: balance})
		} else if balance < 0 {
			heap.Push(&negHeap, Node{UserID: uid, FinalBalance: -balance})
		}
	}

	var graph []PaymentNode

	for posHeap.Len() > 0 && negHeap.Len() > 0 {
		receiver := heap.Pop(&posHeap).(Node)
		sender := heap.Pop(&negHeap).(Node)

		amountTransferred := min(receiver.FinalBalance, sender.FinalBalance)

		graph = append(graph, PaymentNode{From: sender.UserID, To: receiver.UserID, Amount: amountTransferred})

		sender.FinalBalance -= amountTransferred
		receiver.FinalBalance -= amountTransferred

		if sender.FinalBalance > 0 {
			heap.Push(&negHeap, sender)
		}

		if receiver.FinalBalance > 0 {
			heap.Push(&posHeap, receiver)
		}
	}

	return graph
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func addExpense(groupID, expenseID, title, imageURI string, balances map[string]Balance) {
	expense := Expense{
		ExpenseID: expenseID,
		IsSettled: false,
		Balances:  balances,
		GroupID:   groupID,
		Title:     title,
		Timestamp: time.Now().Unix(),
		ImageURI:  imageURI,
	}
	expenses[expenseID] = expense
}

func editExpense(expenseID, title, imageURI string, balances map[string]Balance) {
	if expense, exists := expenses[expenseID]; exists {
		expense.Title = title
		expense.ImageURI = imageURI
		expense.Balances = balances
		expenses[expenseID] = expense
	} else {
		fmt.Println("Expense not found")
	}
}

func settleExpense(expenseID string) {
	if expense, exists := expenses[expenseID]; exists {
		expense.IsSettled = true
		expenses[expenseID] = expense
	} else {
		fmt.Println("Expense not found")
	}
}

func makeGroup(groupID, title, imageURI, description string, users []User) {
	group := Group{
		GroupID:     groupID,
		Users:       users,
		ImageURI:    imageURI,
		Title:       title,
		Description: description,
	}
	groups[groupID] = group
}

func getGroupExpenses(groupID string) []Expense {
	var groupExpenses []Expense
	for _, expense := range expenses {
		if expense.GroupID == groupID {
			groupExpenses = append(groupExpenses, expense)
		}
	}
	return groupExpenses
}

func getGroupPaymentGraph(groupID string) []PaymentNode {
	balances := make(map[string]int)
	for _, expense := range getGroupExpenses(groupID) {
		for userID, balance := range expense.Balances {
			balances[userID] += balance.Amount
		}
	}
	return makePaymentGraph(balances)
}

func getUser(userID string) User {
	return users[userID]
}

func getUsersInGroup(groupID string) []User {
	return groups[groupID].Users
}

func getGroup(groupID string) Group {
	return groups[groupID]
}

func getExpense(expenseID string) Expense {
	return expenses[expenseID]
}

func main() {
	// Sample data
	userA := User{UserID: "A", ImageURI: "", Bio: "User A"}
	userB := User{UserID: "B", ImageURI: "", Bio: "User B"}
	userC := User{UserID: "C", ImageURI: "", Bio: "User C"}
	userD := User{UserID: "D", ImageURI: "", Bio: "User D"}
	userE := User{UserID: "E", ImageURI: "", Bio: "User E"}
	userF := User{UserID: "F", ImageURI: "", Bio: "User F"}
	userG := User{UserID: "G", ImageURI: "", Bio: "User G"}

	users[userA.UserID] = userA
	users[userB.UserID] = userB
	users[userC.UserID] = userC
	users[userD.UserID] = userD
	users[userE.UserID] = userE
	users[userF.UserID] = userF
	users[userG.UserID] = userG

	makeGroup("group1", "Group 1", "", "A sample group", []User{userA, userB, userC, userD, userE, userF, userG})

	balances := map[string]Balance{
		userA.UserID: {Currency: "USD", Amount: 80},
		userB.UserID: {Currency: "USD", Amount: 25},
		userC.UserID: {Currency: "USD", Amount: -25},
		userD.UserID: {Currency: "USD", Amount: -20},
		userE.UserID: {Currency: "USD", Amount: -20},
		userF.UserID: {Currency: "USD", Amount: -20},
		userG.UserID: {Currency: "USD", Amount: -20},
	}

	addExpense("group1", "expense1", "Group Expense", "", balances)

	fmt.Println("Group Expenses:")
	for _, expense := range getGroupExpenses("group1") {
		fmt.Printf("%+v\n", expense)
	}

	paymentGraph := getGroupPaymentGraph("group1")
	fmt.Println("Payment Graph:")
	for _, payment := range paymentGraph {
		fmt.Printf("From: %s, To: %s, Amount: %d\n", payment.From, payment.To, payment.Amount)
	}
}
