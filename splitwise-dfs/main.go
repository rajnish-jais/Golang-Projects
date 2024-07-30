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
	"math"
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

type ExpenseSystem struct {
	Users    map[string]User
	Expenses []Expense
	Balances map[string]map[string]float64
}

type Transaction struct {
	from, to string
	amount   float64
}

func NewExpenseSystem() *ExpenseSystem {
	return &ExpenseSystem{
		Users:    make(map[string]User),
		Balances: make(map[string]map[string]float64),
	}
}

func (es *ExpenseSystem) computeMemberBalances() map[string]float64 {
	memberVsBalanceMap := make(map[string]float64)

	for _, balances := range es.Balances {
		for to, amount := range balances {
			memberVsBalanceMap[to] += amount
		}
	}
	return memberVsBalanceMap
}

func (es *ExpenseSystem) minTransfers() int {
	memberVsBalanceMap := es.computeMemberBalances()

	balanceList := []float64{}
	members := []string{}
	for member, amount := range memberVsBalanceMap {
		if amount != 0 {
			balanceList = append(balanceList, amount)
			members = append(members, member)
		}
	}

	minTxnCount := dfs(balanceList, members, 0)

	return minTxnCount
}

// Helper function for DFS and backtracking
func dfs(balanceList []float64, members []string, currentIndex int) int {
	if len(balanceList) == 0 || currentIndex >= len(balanceList) {
		return 0
	}
	if balanceList[currentIndex] == 0 {
		return dfs(balanceList, members, currentIndex+1)
	}

	currentVal := balanceList[currentIndex]
	minTxnCount := math.MaxInt32

	for txnIndex := currentIndex + 1; txnIndex < len(balanceList); txnIndex++ {
		nextVal := balanceList[txnIndex]
		if currentVal*nextVal < 0 {
			balanceList[txnIndex] += currentVal

			txnCount := 1 + dfs(balanceList, members, currentIndex+1)

			if txnCount < minTxnCount {
				minTxnCount = txnCount
			}

			balanceList[txnIndex] -= currentVal

			if currentVal+nextVal == 0 {
				break
			}
		}
	}

	return minTxnCount
}

func (es *ExpenseSystem) AddUser(name string) {
	es.Users[name] = User{Name: name}
	if es.Balances[name] == nil {
		es.Balances[name] = make(map[string]float64)
	}
}

func (es *ExpenseSystem) AddExpense(payer string, amount float64, participants []string, shareType string, shares map[string]float64, metadata map[string]string) {
	// Compute shares
	var expenseShares []Share
	if shareType == "equal" {
		equalShare := amount / float64(len(participants)+1)
		for _, participant := range participants {
			expenseShares = append(expenseShares, Share{User: participant, Amount: equalShare})
		}
	} else if shareType == "exact" {
		for user, share := range shares {
			expenseShares = append(expenseShares, Share{User: user, Amount: share})
		}
	} else if shareType == "percentage" {
		for user, percentage := range shares {
			shareAmount := (percentage / 100) * amount
			expenseShares = append(expenseShares, Share{User: user, Amount: shareAmount})
		}
	}

	// Add expense
	es.Expenses = append(es.Expenses, Expense{payer, amount, participants, expenseShares, metadata})

	// Update balances
	for _, share := range expenseShares {
		if es.Balances[share.User] == nil {
			es.Balances[share.User] = make(map[string]float64)
		}
		if es.Balances[payer] == nil {
			es.Balances[payer] = make(map[string]float64)
		}
		es.Balances[share.User][payer] += share.Amount
		es.Balances[payer][share.User] -= share.Amount
	}

}

func (es *ExpenseSystem) ListUserExpenses(userName string) []Expense {
	fmt.Printf("User %v's Expenses:\n", userName)
	var userExpenses []Expense
	for _, expense := range es.Expenses {
		if expense.Payer == userName || contains(expense.Participants, userName) {
			userExpenses = append(userExpenses, expense)
		}
	}
	return userExpenses
}

func (es *ExpenseSystem) GenerateIndividualSummary(userName string) {
	fmt.Printf("User %v's Summary:\n", userName)
	for user, balance := range es.Balances[userName] {
		if balance > 0 {
			fmt.Printf("%s owes %s %.2f\n", user, userName, balance)
		} else if balance < 0 {
			fmt.Printf("%s is owed by %s %.2f\n", userName, user, -balance)
		}
	}
}

func (es *ExpenseSystem) SettleFunds(user1, user2 string) {
	fmt.Printf("Settle Funds between %v and %v:\n", user1, user2)
	fmt.Print(es.Balances)
	es.Balances[user1][user2] = 0
	es.Balances[user2][user1] = 0
}

func (es *ExpenseSystem) GenerateOverallSummary() {
	fmt.Println("Overall Summary:")
	for user1, user1Balances := range es.Balances {
		for user2, balance := range user1Balances {
			if balance > 0 {
				fmt.Printf("%s owes %s %.2f\n", user2, user1, balance)
			}
		}
	}
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
	es := NewExpenseSystem()

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

	es.GenerateIndividualSummary("A")
	es.GenerateIndividualSummary("B")
	es.GenerateIndividualSummary("C")
	es.GenerateIndividualSummary("D")

	es.GenerateOverallSummary()

	es.SettleFunds("D", "C")

	es.GenerateIndividualSummary("A")
	es.GenerateIndividualSummary("B")
	es.GenerateIndividualSummary("C")
	es.GenerateIndividualSummary("D")

	minTxnCount := es.minTransfers()

	fmt.Println("Minimum number of transactions required:", minTxnCount)
	es.GenerateOverallSummary()
}
