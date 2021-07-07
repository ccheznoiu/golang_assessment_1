package main

type plansResponse struct {
	//ID              int     `json:"id"`
	//DebtID          int     `json:"debt_id"`
	AmtToPay        USD    `json:"amount_to_pay"`
	InstallmentFreq string `json:"installment_frequency"`
	//InstallmentAmt  float64 `json:"installment_amount"`
	StartDate string `json:"start_date"`
}
