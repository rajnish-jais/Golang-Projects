package models

type Expense struct {
	Payer        string
	Amount       float64
	Participants []string
	Shares       []Share
	Metadata     map[string]string
}
