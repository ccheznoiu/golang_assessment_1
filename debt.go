package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type debtResponse struct {
	ID           *int    `json:"id"`
	Amount       USD     `json:"amount"`
	IsInPayPlan  bool    `json:"is_in_payment_plan"`
	RemainingAmt USD     `json:"remaining_amount"`
	NextPayDue   *string `json:"next_payment_due_date"`
}

func (d *debtResponse) leveragePlan() {
	if d.ID == nil {
		return
	}

	plansBody := get(plansURL, map[string]int{"debt_id": *d.ID})
	defer plansBody.Close()

	var pf int

	plans := decode(plansBody)
	for plans.More() {
		var plan plansResponse
		if err := plans.Decode(&plan); err != nil {
			log.Fatal(err)
		}

		d.RemainingAmt = plan.AmtToPay
		d.NextPayDue = &plan.StartDate
		if pf = freqMap[plan.InstallmentFreq]; pf == 0 {
			log.Printf("encountered an invalid installment_frequency: %s!\n", plan.InstallmentFreq)
		}
	}

	d.creditPayments(pf)
}

var freqMap = map[string]int{
	"WEEKLY":    7,
	"BI_WEEKLY": 14,
}

func (d *debtResponse) creditPayments(pf int) {
	if d.ID == nil {
		return
	}

	paymentsBody := get(paymentsURL, map[string]int{"payment_plan_id": *d.ID})
	defer paymentsBody.Close()

	var nextInstall time.Time
	var err error

	if d.NextPayDue != nil {
		if nextInstall, err = time.Parse(ISO8601, *d.NextPayDue); err != nil {
			log.Printf("encountered invalid start_date: %s!\n", *d.NextPayDue)
			d.NextPayDue = nil
		} else if time.Now().Truncate(24 * time.Hour).After(nextInstall) {
			d.IsInPayPlan = true
		}
	}

	payments := decode(paymentsBody)
	for payments.More() {
		var payment paymentResponse
		if err = payments.Decode(&payment); err != nil {
			log.Fatal(err)
		}

		d.RemainingAmt -= payment.Amount

		pDate, err := time.Parse(ISO8601, payment.Date)
		if err != nil {
			log.Printf("encountered invalid payment date: %s!\n", payment.Date)
			continue
		}

		for nextInstall.Before(pDate) {
			nextInstall = nextInstall.AddDate(0, 0, pf)
		}
	}

	if d.RemainingAmt <= 0 {
		d.IsInPayPlan = false
		d.NextPayDue = nil

		return
	}

	if d.NextPayDue != nil {
		nextPayDue := nextInstall.Format(ISO8601)
		d.NextPayDue = &nextPayDue
	}
}

func (d debtResponse) printOut() {
	debtOut, err := json.MarshalIndent(d, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(debtOut))
}
