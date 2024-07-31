package services

import (
	"fmt"
	"math"
	"sync"

	"splitwise/models"
	"splitwise/utils"
)

//type ExpenseSystem interface {
//	AddUser(name string)
//	AddExpense(payer string, amount float64, participants []string, shareType string, shares map[string]float64, metadata map[string]string)
//	ListUserExpenses(userName string) []models.Expense
//	GenerateIndividualSummary(userName string) []string
//	SettleFunds(user1, user2 string)
//	GenerateOverallSummary() []string
//}

type SimpleExpenseSystem struct {
	mu       sync.RWMutex
	users    map[string]models.User
	expenses []models.Expense
	balances map[string]map[string]float64
}

func NewSimpleExpenseSystem() *SimpleExpenseSystem {
	return &SimpleExpenseSystem{
		users:    make(map[string]models.User),
		balances: make(map[string]map[string]float64),
	}
}

func (es *SimpleExpenseSystem) AddUser(name string) {
	es.mu.Lock()
	defer es.mu.Unlock()
	es.users[name] = models.User{Name: name}
	if es.balances[name] == nil {
		es.balances[name] = make(map[string]float64)
	}
}

func (es *SimpleExpenseSystem) AddExpense(payer string, amount float64, participants []string, shareType string, shares map[string]float64, metadata map[string]string) {
	var expenseShares []models.Share
	switch shareType {
	case "equal":
		equalShare := amount / float64(len(participants)+1)
		for _, participant := range participants {
			expenseShares = append(expenseShares, models.Share{User: participant, Amount: equalShare})
		}
	case "exact":
		for user, share := range shares {
			expenseShares = append(expenseShares, models.Share{User: user, Amount: share})
		}
	case "percentage":
		for user, percentage := range shares {
			shareAmount := (percentage / 100) * amount
			expenseShares = append(expenseShares, models.Share{User: user, Amount: shareAmount})
		}
	}

	es.mu.Lock()
	es.expenses = append(es.expenses, models.Expense{payer, amount, participants, expenseShares, metadata})

	for _, share := range expenseShares {
		es.balances[share.User][payer] += share.Amount
		es.balances[payer][share.User] -= share.Amount
	}
	es.mu.Unlock()
}

func (es *SimpleExpenseSystem) ListUserExpenses(userName string) []models.Expense {
	es.mu.RLock()
	defer es.mu.RUnlock()
	fmt.Printf("\nUser %v's Expenses:\n", userName)

	var userExpenses []models.Expense
	for _, expense := range es.expenses {
		if expense.Payer == userName || utils.Contains(expense.Participants, userName) {
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

func (es *SimpleExpenseSystem) computeMemberBalances() map[string]float64 {
	memberVsBalanceMap := make(map[string]float64)

	for _, balances := range es.balances {
		for to, amount := range balances {
			memberVsBalanceMap[to] += amount
		}
	}
	return memberVsBalanceMap
}

func (es *SimpleExpenseSystem) MinTransfers() int {
	memberVsBalanceMap := es.computeMemberBalances()

	balanceList := []float64{}
	for _, amount := range memberVsBalanceMap {
		if amount != 0 {
			balanceList = append(balanceList, amount)
		}
	}

	minTxnCount := dfs(balanceList, 0)

	return minTxnCount
}

// Helper function for DFS and backtracking
func dfs(balanceList []float64, currentIndex int) int {
	if len(balanceList) == 0 || currentIndex >= len(balanceList) {
		return 0
	}
	if balanceList[currentIndex] == 0 {
		return dfs(balanceList, currentIndex+1)
	}

	currentVal := balanceList[currentIndex]
	minTxnCount := math.MaxInt32

	for txnIndex := currentIndex + 1; txnIndex < len(balanceList); txnIndex++ {
		nextVal := balanceList[txnIndex]
		if currentVal*nextVal < 0 {
			balanceList[txnIndex] += currentVal

			txnCount := 1 + dfs(balanceList, currentIndex+1)

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
