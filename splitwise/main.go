/*
Design an Expense Sharing System

You are tasked with designing an expense sharing system that allows users to track shared expenses among a group of people. The system should support the following functionalities:

Add Expense: Implement a function that allows a user to add an expense. The function should take the following parameters:

Main user (the one who paid)
Amount of the expense
List of users included in the payment
Shares of each user (which can be equal, exact amounts, or percentages)
Metadata about the expense
List User's Expenses: Create a function that lists all the expenses a given user is involved in. This should include details of the expenses, such as the amount, shares, and the other participants.

Generate Individual Summary: Implement a function that generates a summary for each individual user. This summary should show all the expenses they are involved in and whether they owe money or are owed money.

Settle Funds: Develop a function that calculates and performs fund settlements between two users. Given two users, find the simplest way to settle the debts, if any exist.

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
)

// Expense represents a single expense
type Expense struct {
	PaidBy    string
	Amount    int
	Users     []string
	ShareType string
	Shares    map[string]int
}

// UserSummary represents a summary for an individual user
type UserSummary struct {
	Owes      map[string]int
	OwedBy    map[string]int
	TotalOwes int
	TotalOwed int
}

// ExpenseSystem represents the expense sharing system
type ExpenseSystem struct {
	Expenses      []Expense
	UserSummaries map[string]UserSummary
	UserExpenses  map[string][]Expense
}

// AddExpense adds a new expense to the system
func (es *ExpenseSystem) AddExpense(paidBy string, amount int, users []string, shareType string, shares map[string]int) {
	expense := Expense{
		PaidBy:    paidBy,
		Amount:    amount,
		Users:     users,
		ShareType: shareType,
		Shares:    shares,
	}
	es.Expenses = append(es.Expenses, expense)

	// Update UserSummaries and UserExpenses
	if _, ok := es.UserSummaries[paidBy]; !ok {
		es.UserSummaries[paidBy] = UserSummary{Owes: make(map[string]int), OwedBy: make(map[string]int)}
	}

	for _, u := range users {
		if u != paidBy {
			if _, ok := es.UserSummaries[u]; !ok {
				es.UserSummaries[u] = UserSummary{Owes: make(map[string]int), OwedBy: make(map[string]int)}
			}

			// Temporarily store the summary
			tempSummary := es.UserSummaries[paidBy]
			tempSummary.TotalOwes += shares[u]
			tempSummary.Owes[u] += shares[u]

			es.UserSummaries[paidBy] = tempSummary

			tempSummary = es.UserSummaries[u]
			tempSummary.TotalOwed += shares[u]
			tempSummary.OwedBy[paidBy] += shares[u]

			es.UserSummaries[u] = tempSummary

			// Update UserExpenses
			es.UserExpenses[u] = append(es.UserExpenses[u], expense)
		}
	}
}

// GenerateIndividualSummary generates a summary for an individual user
func (es *ExpenseSystem) GenerateIndividualSummary(user string) UserSummary {
	return es.UserSummaries[user]
}

// ListUserExpenses lists all the expenses a given user is involved in
func (es *ExpenseSystem) ListUserExpenses(user string) []Expense {
	return es.UserExpenses[user]
}

func main() {
	expenseSystem := ExpenseSystem{
		UserSummaries: make(map[string]UserSummary),
		UserExpenses:  make(map[string][]Expense),
	}

	// Example 1
	expenseSystem.AddExpense("A", 1000, []string{"B", "C", "D"}, "equal", nil)

	// Example 2
	expenseSystem.AddExpense("B", 1000, []string{"A", "C"}, "exact", map[string]int{"A": 300, "C": 500})

	// Example 3
	expenseSystem.AddExpense("C", 1000, []string{"A", "B", "D"}, "percentage", map[string]int{"A": 25, "B": 30, "D": 10})

	// User A's Summary
	userASummary := expenseSystem.GenerateIndividualSummary("A")
	fmt.Println("User A's Summary:")
	fmt.Println("Total Owes:", userASummary.TotalOwes)
	fmt.Println("Details:")
	for user, amount := range userASummary.Owes {
		fmt.Printf("%s owes %s: %d\n", user, "A", amount)
	}

	// User B's Summary
	userBSummary := expenseSystem.GenerateIndividualSummary("B")
	fmt.Println("\nUser B's Summary:")
	fmt.Println("Total Owes:", userBSummary.TotalOwes)
	fmt.Println("Details:")
	for user, amount := range userBSummary.Owes {
		fmt.Printf("%s owes %s: %d\n", user, "B", amount)
	}

	// Listing User A's Expenses
	userAExpenses := expenseSystem.ListUserExpenses("A")
	fmt.Println("\nUser A's Expenses:")
	for i, expense := range userAExpenses {
		fmt.Printf("%d. %s Pays %d with %v (%s)\n", i+1, expense.PaidBy, expense.Amount, expense.Users, expense.ShareType)
	}
}
