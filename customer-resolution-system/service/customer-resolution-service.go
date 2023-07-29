package service

import "fmt"

// IssueType represents the type of customer issue
type IssueType string

const (
	PaymentRelated    = "Payment Related"
	MutualFundRelated = "Mutual Fund Related"
	GoldRelated       = "Gold Related"
	InsuranceRelated  = "Insurance Related"
)

// IssueStatus represents the status of a customer issue
type IssueStatus string

const (
	Open       = "Open"
	Assigned   = "Assigned"
	InProgress = "In Progress"
	Resolved   = "Resolved"
)

// Issue represents a customer issue
type Issue struct {
	ID          string
	Transaction string
	IssueType   IssueType
	Subject     string
	Description string
	Email       string
	Status      IssueStatus
	Resolution  string
}

// Agent represents a customer service agent
type Agent struct {
	Email     string
	Name      string
	Expertise []string
}

// CustomerService represents the main system
type CustomerService struct {
	issues       map[string]Issue
	agents       map[string]Agent
	assigned     map[string]bool
	waitingQueue []string
}

// Create a new customer service
func NewCustomerService() *CustomerService {
	return &CustomerService{
		issues:       make(map[string]Issue),
		agents:       make(map[string]Agent),
		assigned:     make(map[string]bool),
		waitingQueue: []string{},
	}
}

// Function to create a new issue
func (cs *CustomerService) CreateIssue(transactionID string, issueType IssueType, subject, description, email string) {
	issueID := "I" + fmt.Sprint(len(cs.issues)+1)

	issue := Issue{
		ID:          issueID,
		Transaction: transactionID,
		IssueType:   issueType,
		Subject:     subject,
		Description: description,
		Email:       email,
		Status:      Open,
		Resolution:  "",
	}

	cs.issues[issueID] = issue
	fmt.Printf(">>> Issue %s created against transaction \"%s\"\n", issueID, transactionID)

	// Assign the issue if an agent is available
	cs.AssignIssue(issueID)
}

// Function to add a new agent
func (cs *CustomerService) AddAgent(agentEmail, agentName string, expertise []string) {
	agent := Agent{
		Email:     agentEmail,
		Name:      agentName,
		Expertise: expertise,
	}

	cs.agents[agentEmail] = agent
	fmt.Printf(">>> Agent %s created\n", agentEmail)
}

// Function to assign an issue to an agent
func (cs *CustomerService) AssignIssue(issueID string) {
	for email, _ := range cs.agents {
		if !cs.assigned[email] {
			issue, ok := cs.issues[issueID]
			if !ok {
				fmt.Println("Issue not found")
				return
			}

			issue.Status = Assigned
			cs.issues[issueID] = issue

			cs.assigned[email] = true
			fmt.Printf(">>> Issue %s assigned to agent %s\n", issueID, email)
			return
		}
	}

	// If no agent is available, add the issue to the waiting queue
	cs.waitingQueue = append(cs.waitingQueue, issueID)
	//fmt.Printf(">>> Issue %s added to the waitlist\n", issueID)
}

// Function to retrieve issues based on a filter (e.g., email or issue type)
func (cs *CustomerService) GetIssues(filter map[string]string) []Issue {
	var filteredIssues []Issue

	for _, issue := range cs.issues {
		match := true
		for key, value := range filter {
			if key == "email" {
				if issue.Email != value {
					match = false
					break
				}
			} else if key == "type" {
				if issue.IssueType != IssueType(value) {
					match = false
					break
				}
			}
		}

		if match {
			filteredIssues = append(filteredIssues, issue)
		}
	}

	return filteredIssues
}

// Function to update an issue's status and resolution
func (cs *CustomerService) UpdateIssue(issueID string, status IssueStatus, resolution string) {
	issue, ok := cs.issues[issueID]
	if !ok {
		fmt.Println("Issue not found")
		return
	}

	issue.Status = status
	issue.Resolution = resolution

	cs.issues[issueID] = issue
	fmt.Printf(">>> Issue %s status updated to %s\n", issueID, status)
}

// Function to resolve an issue
func (cs *CustomerService) ResolveIssue(issueID, resolution string) {
	issue, ok := cs.issues[issueID]
	if !ok {
		fmt.Println("Issue not found")
		return
	}

	issue.Status = Resolved
	issue.Resolution = resolution

	cs.issues[issueID] = issue
	fmt.Printf(">>> Issue %s marked as resolved\n", issueID)

	// After resolving the issue, check if there are waiting issues for agents
	if len(cs.waitingQueue) > 0 {
		cs.AssignIssue(cs.waitingQueue[0])
		cs.waitingQueue = cs.waitingQueue[1:]
	}
}

// Function to view agents' work history
func (cs *CustomerService) ViewAgentsWorkHistory() map[string][]Issue {
	agentWorkHistory := make(map[string][]Issue)

	for email, _ := range cs.agents {
		agentIssues := []Issue{}
		for _, issue := range cs.issues {
			if cs.assigned[email] && issue.Status == Assigned {
				agentIssues = append(agentIssues, issue)
			}
		}
		agentWorkHistory[email] = agentIssues
	}

	return agentWorkHistory
}
