package main

import "log"

var (
	debtsURL    string = "/debts"
	plansURL    string = "/payment_plans"
	paymentsURL string = "/payments"
)

const ISO8601 string = "2006-01-02"

func main() {
	debtsBody := get(debtsURL, nil)
	defer debtsBody.Close()

	debts := decode(debtsBody)
	for debts.More() {
		var debt debtResponse
		if err := debts.Decode(&debt); err != nil {
			log.Fatal(err)
		} else if debt.ID == nil {
			log.Println("encountered invalid debt object!")
			continue
		}

		debt.RemainingAmt = debt.Amount
		debt.leveragePlan()
		debt.printOut()
	}
}
