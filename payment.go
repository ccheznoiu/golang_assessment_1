package main

type paymentResponse struct {
	//PaymentPlanID int    `json:"payment_plan_id"`
	Amount USD    `json:"amount"`
	Date   string `json:"date"`
}
