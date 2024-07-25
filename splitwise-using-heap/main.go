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
	"errors"
	"fmt"
)

// User struct definition
type User struct {
	UserID   string
	ImageURI string
	Bio      string
}

// Balance struct definition
type Balance struct {
	Currency string
	Amount   int
}

// Expense struct definition
type Expense struct {
	ExpenseID    string
	IsSettled    bool
	UserBalances map[string]Balance // maps userID to Balance
	GroupID      string
	Title        string
	Timestamp    int64
	ImageURI     string
}

// Group struct definition
type Group struct {
	GroupID     string
	Users       []User
	ImageURI    string
	Title       string
	Description string
}

// PaymentNode struct definition for balancing algorithm
type PaymentNode struct {
	From   string
	To     string
	Amount int
}

// MinHeap for balancing algorithm
type MinHeap []Node

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i].FinalBalance < h[j].FinalBalance }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MinHeap) Push(x interface{}) {
	*h = append(*h, x.(Node))
}

func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// Node struct for balancing algorithm
type Node struct {
	UserID       string
	FinalBalance int
}

// In-memory database for users, groups, and expenses
var users = map[string]User{}
var groups = map[string]Group{}
var expenses = map[string]Expense{}

// Helper function to get the group by ID
func getGroup(groupID string) (Group, error) {
	group, exists := groups[groupID]
	if !exists {
		return Group{}, errors.New("group not found")
	}
	return group, nil
}

// Helper function to get the user by ID
func getUser(userID string) (User, error) {
	user, exists := users[userID]
	if !exists {
		return User{}, errors.New("user not found")
	}
	return user, nil
}

// AddExpense function
func addExpense(groupID string, userID string, amount int, currency string, title string, timestamp int64, imageURI string) string {
	expenseID := fmt.Sprintf("%s_%d", groupID, len(expenses)+1)
	userBalances := map[string]Balance{
		userID: {Currency: currency, Amount: amount},
	}
	expenses[expenseID] = Expense{
		ExpenseID:    expenseID,
		IsSettled:    false,
		UserBalances: userBalances,
		GroupID:      groupID,
		Title:        title,
		Timestamp:    timestamp,
		ImageURI:     imageURI,
	}
	return expenseID
}

// EditExpense function
func editExpense(expenseID string, userID string, amount int, currency string) error {
	expense, exists := expenses[expenseID]
	if !exists {
		return errors.New("expense not found")
	}
	expense.UserBalances[userID] = Balance{Currency: currency, Amount: amount}
	expenses[expenseID] = expense
	return nil
}

// SettleExpense function
func settleExpense(expenseID string) error {
	expense, exists := expenses[expenseID]
	if !exists {
		return errors.New("expense not found")
	}
	expense.IsSettled = true
	expenses[expenseID] = expense
	return nil
}

// CreateGroup function
func createGroup(users []User, title string, imageURI string, description string) string {
	groupID := fmt.Sprintf("group_%d", len(groups)+1)
	groups[groupID] = Group{
		GroupID:     groupID,
		Users:       users,
		ImageURI:    imageURI,
		Title:       title,
		Description: description,
	}
	return groupID
}

// GetGroupExpenses function
func getGroupExpenses(groupID string) ([]Expense, error) {
	var groupExpenses []Expense
	for _, expense := range expenses {
		if expense.GroupID == groupID {
			groupExpenses = append(groupExpenses, expense)
		}
	}
	return groupExpenses, nil
}

// makePaymentGraph function
func makePaymentGraph(groupID string) ([]PaymentNode, error) {
	_, err := getGroup(groupID)
	if err != nil {
		return nil, err
	}

	// Calculate final balances
	finalBalances := map[string]int{}
	for _, expense := range expenses {
		if expense.GroupID == groupID {
			for userID, balance := range expense.UserBalances {
				finalBalances[userID] += balance.Amount
			}
		}
	}

	positiveBalances := &MinHeap{}
	negativeBalances := &MinHeap{}
	heap.Init(positiveBalances)
	heap.Init(negativeBalances)

	// Divide users into positive and negative balance heaps
	for userID, balance := range finalBalances {
		if balance > 0 {
			heap.Push(positiveBalances, Node{UserID: userID, FinalBalance: balance})
		} else if balance < 0 {
			heap.Push(negativeBalances, Node{UserID: userID, FinalBalance: -balance})
		}
	}

	var paymentGraph []PaymentNode

	// Balancing algorithm
	for positiveBalances.Len() > 0 && negativeBalances.Len() > 0 {
		receiver := heap.Pop(positiveBalances).(Node)
		sender := heap.Pop(negativeBalances).(Node)

		amountTransferred := min(receiver.FinalBalance, sender.FinalBalance)
		paymentGraph = append(paymentGraph, PaymentNode{From: sender.UserID, To: receiver.UserID, Amount: amountTransferred})

		receiver.FinalBalance -= amountTransferred
		sender.FinalBalance -= amountTransferred

		if receiver.FinalBalance != 0 {
			heap.Push(positiveBalances, receiver)
		}
		if sender.FinalBalance != 0 {
			heap.Push(negativeBalances, sender)
		}
	}

	return paymentGraph, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	// Example usage
	u1 := User{UserID: "u1", ImageURI: "img1.jpg", Bio: "User 1"}
	u2 := User{UserID: "u2", ImageURI: "img2.jpg", Bio: "User 2"}
	u3 := User{UserID: "u3", ImageURI: "img3.jpg", Bio: "User 3"}

	users[u1.UserID] = u1
	users[u2.UserID] = u2
	users[u3.UserID] = u3

	groupID := createGroup([]User{u1, u2, u3}, "Group 1", "group_img.jpg", "Group Description")
	addExpense(groupID, u1.UserID, 100, "USD", "Lunch", 1625154000, "lunch_img.jpg")
	addExpense(groupID, u2.UserID, -50, "USD", "Dinner", 1625154001, "dinner_img.jpg")
	addExpense(groupID, u3.UserID, -50, "USD", "Snacks", 1625154002, "snacks_img.jpg")

	paymentGraph, _ := makePaymentGraph(groupID)
	for _, payment := range paymentGraph {
		fmt.Printf("User %s pays %d to User %s\n", payment.From, payment.Amount, payment.To)
	}
}
