package main

/*
An individual payment installment which is made on a payment plan. Many-to-one with debts.

payment_plan_id (integer)
amount (real)
date (string) - ISO 8601 date of when this payment occurred.
*/

type paymentResponse struct {
	//PaymentPlanID int    `json:"payment_plan_id"`
	Amount USD    `json:"amount"`
	Date   string `json:"date"`
}
