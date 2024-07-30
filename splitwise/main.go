/*
Design an Expense Sharing System

You are tasked with designing an expense sharing system that allows users to track shared expenses among a group of people.
The system should support the following functionalities:

Settle Funds: Develop a function that calculates and performs fund settlements between two users. Given two users, find the simplest way to settle the debts, if any exist.

Add Expense: Implement a function that allows a user to add an expense. The function should take the following parameters:

Main user (the one who paid)
Amount of the expense
List of users included in the payment
Shares of each user (which can be equal, exact amounts, or percentages)
Metadata about the expense

List User's Expenses: Create a function that lists all the expenses a given user is involved in. This should include details of the expenses, such as the amount, shares, and the other participants.

Generate Individual Summary: Implement a function that generates a summary for each individual user. This summary should show all the expenses they are involved in and whether they owe money or are owed money.

Generate Overall Summary: Create a function that generates an overall summary of the expenses. The summary should include details of who owes whom and how much.

Use the following examples as a guide to test your system:

User A pays 1000 for a payment with users B, C, and D included. All shares are equal.
User B pays 1000 for a payment with users A and C included. A shares 300, C shares 500.
User C pays 1000 for a payment with users A, B, and D included. A percentage 25, B percentage 30, D percentage 10.
# Example 1: A pays 1000 for a payment with B, C, D included -> All shares are equal (1000 / 4)
expense_system.add_expense("A", 1000, ["B", "C", "D"], "equal")
# Expected Summary:
# B owes A 250
# C owes A 250
# D owes A 250

# Example 2: B pays 1000 for a payment with A, C included -> A shares 300, C shares 500
expense_system.add_expense("B", 1000, ["A", "C"], "exact", {"A": 300, "C": 500})
# Expected Summary:
# A owes B 50 (300 - 250)
# C owes A 250
# C owes B 500
# D owes A 250

# Example 3: C pays 1000 for a payment with A, B, D included -> A percentage 25, B percentage 30, D percentage 10
expense_system.add_expense("C", 1000, ["A", "B", "D"], "percentage", {"A": 25, "B": 30, "D": 10})
# Expected Summary:
# A owes B 50 (300 - 250)
# C owes B 200
# D owes A 250
# D owes C 100

# User A's Summary
user_a_summary = expense_system.generate_individual_summary("A")
# Expected User A's Summary:
# A owes B 50
# C owes A 250
# D owes A 250

# User B's Summary
user_b_summary = expense_system.generate_individual_summary("B")
# Expected User B's Summary:
# B owes A 250

# User C's Summary
user_c_summary = expense_system.generate_individual_summary("C")
# Expected User C's Summary:
# C owes A 250

# User D's Summary
user_d_summary = expense_system.generate_individual_summary("D")
# Expected User D's Summary:
# D owes A 250
# D owes C 100

# Listing User A's Expenses
user_a_expenses = expense_system.list_user_expenses("A")
# Expected User A's Expenses:
# 1. A Pays 1000 with B, C, and D (EQUAL)
# 2. B Pays 1000 with A and C (EXACT => A = 300 C = 500)
# 3. C Pays 1000 with A, B, and D (PERCENTAGE => A = 25 B = 30 D = 10)
*/
package main

import (
	"fmt"
	"splitwise/services"
	"sync"
)

type User struct {
	Name string
}

type Share struct {
	User   string
	Amount float64
}

type Expense struct {
	Payer        string
	Amount       float64
	Participants []string
	Shares       []Share
	Metadata     map[string]string
}

type SimpleExpenseSystem struct {
	mu       sync.RWMutex
	users    map[string]User
	expenses []Expense
	balances map[string]map[string]float64
}

func NewSimpleExpenseSystem() *SimpleExpenseSystem {
	return &SimpleExpenseSystem{
		users:    make(map[string]User),
		balances: make(map[string]map[string]float64),
	}
}

func (es *SimpleExpenseSystem) AddUser(name string) {
	es.mu.Lock()
	defer es.mu.Unlock()
	es.users[name] = User{Name: name}
	if es.balances[name] == nil {
		es.balances[name] = make(map[string]float64)
	}
}

func (es *SimpleExpenseSystem) AddExpense(payer string, amount float64, participants []string, shareType string, shares map[string]float64, metadata map[string]string) {
	var expenseShares []Share
	switch shareType {
	case "equal":
		equalShare := amount / float64(len(participants)+1)
		for _, participant := range participants {
			expenseShares = append(expenseShares, Share{User: participant, Amount: equalShare})
		}
	case "exact":
		for user, share := range shares {
			expenseShares = append(expenseShares, Share{User: user, Amount: share})
		}
	case "percentage":
		for user, percentage := range shares {
			shareAmount := (percentage / 100) * amount
			expenseShares = append(expenseShares, Share{User: user, Amount: shareAmount})
		}
	}

	es.mu.Lock()
	es.expenses = append(es.expenses, Expense{payer, amount, participants, expenseShares, metadata})

	for _, share := range expenseShares {
		es.balances[share.User][payer] += share.Amount
		es.balances[payer][share.User] -= share.Amount
	}
	es.mu.Unlock()
}

func (es *SimpleExpenseSystem) ListUserExpenses(userName string) []Expense {
	es.mu.RLock()
	defer es.mu.RUnlock()
	fmt.Printf("\nUser %v's Expenses:\n", userName)

	var userExpenses []Expense
	for _, expense := range es.expenses {
		if expense.Payer == userName || contains(expense.Participants, userName) {
			userExpenses = append(userExpenses, expense)
		}
	}
	return userExpenses
}

func (es *SimpleExpenseSystem) GenerateIndividualSummary(userName string) []string {
	es.mu.RLock()
	defer es.mu.RUnlock()
	fmt.Printf("\nUser %v's Summary:\n", userName)

	var summary []string
	for user, balance := range es.balances[userName] {
		if balance > 0 {
			summary = append(summary, fmt.Sprintf("%s owes %s %.2f", user, userName, balance))
		} else if balance < 0 {
			summary = append(summary, fmt.Sprintf("%s is owed by %s %.2f", userName, user, -balance))
		}
	}
	return summary
}

func (es *SimpleExpenseSystem) SettleFunds(user1, user2 string) {
	es.mu.Lock()
	defer es.mu.Unlock()
	fmt.Printf("\nSettle Funds between %v and %v:\n", user1, user2)

	es.balances[user1][user2] = 0
	es.balances[user2][user1] = 0
}

func (es *SimpleExpenseSystem) GenerateOverallSummary() []string {
	es.mu.RLock()
	defer es.mu.RUnlock()
	fmt.Println("\nOverall Summary:")
	var summary []string
	for user1, user1Balances := range es.balances {
		for user2, balance := range user1Balances {
			if balance > 0 {
				summary = append(summary, fmt.Sprintf("%s owes %s %.2f", user2, user1, balance))
			}
		}
	}
	return summary
}

func contains(slice []string, item string) bool {
	for _, element := range slice {
		if element == item {
			return true
		}
	}
	return false
}

func main() {
	//es := NewSimpleExpenseSystem()
	es := services.NewSimpleExpenseSystem()

	es.AddUser("A")
	es.AddUser("B")
	es.AddUser("C")
	es.AddUser("D")

	es.AddExpense("A", 1000, []string{"B", "C", "D"}, "equal", nil, nil)
	es.AddExpense("B", 1000, []string{"A", "C"}, "exact", map[string]float64{"A": 300, "C": 500}, nil)
	es.AddExpense("C", 1000, []string{"A", "B", "D"}, "percentage", map[string]float64{"A": 25, "B": 30, "D": 10}, nil)

	for _, expense := range es.ListUserExpenses("A") {
		fmt.Printf("%+v\n", expense)
	}

	for _, expense := range es.ListUserExpenses("D") {
		fmt.Printf("%+v\n", expense)
	}

	for _, summary := range es.GenerateIndividualSummary("A") {
		fmt.Println(summary)
	}

	for _, summary := range es.GenerateIndividualSummary("B") {
		fmt.Println(summary)
	}

	for _, summary := range es.GenerateIndividualSummary("C") {
		fmt.Println(summary)
	}

	for _, summary := range es.GenerateIndividualSummary("D") {
		fmt.Println(summary)
	}

	for _, summary := range es.GenerateOverallSummary() {
		fmt.Println(summary)
	}

	fmt.Println("Minimum number of transactions required:", es.MinTransfers())

	es.SettleFunds("D", "C")

	fmt.Println("Minimum number of transactions required:", es.MinTransfers())

	for _, summary := range es.GenerateIndividualSummary("A") {
		fmt.Println(summary)
	}

	for _, summary := range es.GenerateIndividualSummary("B") {
		fmt.Println(summary)
	}

	for _, summary := range es.GenerateIndividualSummary("C") {
		fmt.Println(summary)
	}

	for _, summary := range es.GenerateIndividualSummary("D") {
		fmt.Println(summary)
	}

	for _, summary := range es.GenerateOverallSummary() {
		fmt.Println(summary)
	}
}
