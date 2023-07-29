package main

import (
	"customer-resolution-service/service"
	"encoding/json"
	"fmt"
)

func main() {
	cs := service.NewCustomerService()

	cs.CreateIssue("T1", service.PaymentRelated, "Payment Failed", "My payment failed but money is debited", "testUser1@test.com")
	cs.CreateIssue("T2", service.MutualFundRelated, "Purchase Failed", "Unable to purchase Mutual Fund", "testUser2@test.com")
	cs.CreateIssue("T3", service.PaymentRelated, "Payment Failed", "My payment failed but money is debited", "testUser2@test.com")

	cs.AddAgent("agent1@test.com", "Agent 1", []string{"Payment Related", "Gold Related"})
	cs.AddAgent("agent2@test.com", "Agent 2", []string{"Payment Related"})

	cs.AssignIssue("I1")
	cs.AssignIssue("I2")
	cs.AssignIssue("I3")

	fmt.Println("Fetching Issues:")
	fmt.Println(cs.GetIssues(map[string]string{"email": "testUser2@test.com"}))

	fmt.Println()
	fmt.Println(cs.GetIssues(map[string]string{"type": service.PaymentRelated}))

	fmt.Println()
	cs.UpdateIssue("I3", service.InProgress, "Waiting for payment confirmation")
	cs.ResolveIssue("I3", "PaymentFailed debited amount will get reversed")

	response, err := json.MarshalIndent(cs.ViewAgentsWorkHistory(), "", "  ")
	if err != nil {
		fmt.Println("Got the error:", err)
	}
	fmt.Println(string(response))
}
