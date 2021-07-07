package main

/*
A payment plan, which is an amount needed to resolve a debt, as well as the frequency of when it will be paid. Payment plans are associated with exactly one debt, and debts may not be associated with more than one payment plan.

id (integer)
debt_id (integer) - The associated debt.
amount_to_pay (real) - Total amount (in USD) needed to be paid to resolve this payment plan.
installment_frequency (text) - The frequency of payments. Is one of: WEEKLY or BI_WEEKLY (14 days).
installment_amount (real) - The amount (in USD) of each payment installment.
start_date (string) - ISO 8601 date of when the first payment is due.
*/

type plansResponse struct {
	//ID              int     `json:"id"`
	//DebtID          int     `json:"debt_id"`
	AmtToPay        USD    `json:"amount_to_pay"`
	InstallmentFreq string `json:"installment_frequency"`
	//InstallmentAmt  float64 `json:"installment_amount"`
	StartDate string `json:"start_date"`
}
