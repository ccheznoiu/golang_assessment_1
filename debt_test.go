package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRemainingAmt(t *testing.T) {
	var ID int
	nextPay := "2020-08-01"
	var testDebt = debtResponse{
		ID:           &ID,
		RemainingAmt: 100, // $1
		IsInPayPlan:  true,
		NextPayDue:   &nextPay,
	}

	fmt.Println("Debt:")
	testDebt.printOut()
	fmt.Println("")

	api := applyPayments(10, 0, 10, time.Time{})
	defer api.Close()

	planAPI := testServer([]byte(`[]`))
	plansURL = planAPI.URL
	defer planAPI.Close()

	testDebt.leveragePlan()
	fmt.Printf("\nDebt after:\n")
	testDebt.printOut()
	fmt.Println("")

	if testDebt.IsInPayPlan || testDebt.RemainingAmt != 0 || testDebt.NextPayDue != nil {
		t.Fail()
	}
}

func TestNextPaymentDue(t *testing.T) {
	var ID int
	testDebt := debtResponse{
		ID:           &ID,
		RemainingAmt: 1000, // $10
	}

	fmt.Println("Debt:")
	testDebt.printOut()
	fmt.Println("")

	sep1 := time.Date(2000, time.September, 1, 0, 0, 0, 0, time.UTC)
	aug1 := time.Date(2000, time.August, 1, 0, 0, 0, 0, time.UTC)

	for k, v := range map[string]int{
		"WEEKLY":    7,
		"BI_WEEKLY": 14,
	} {
		testPlan := plansResponse{
			AmtToPay:        100000, // $1000.00
			InstallmentFreq: k,
			StartDate:       "2000-08-01",
		}

		api := applyPlan(testPlan)
		defer api.Close()

		payAPI := applyPayments(1, 0, 1, sep1)
		defer payAPI.Close()

		testDebt.leveragePlan()
		fmt.Printf("\nDebt after:\n")
		testDebt.printOut()
		fmt.Println("")

		dur := int(sep1.Sub(aug1).Hours()) / 24
		cycles := dur / v
		if dur%v != 0 {
			cycles++
		}

		if *testDebt.NextPayDue != aug1.AddDate(0, 0, v*cycles).Format(ISO8601) {
			t.Fail()
		}
	}
}

func applyPlan(test plansResponse) *httptest.Server {
	planStream := []byte(`[`)
	planBytes, _ := json.Marshal(test)
	planStream = append(planStream, planBytes...)
	planStream = append(planStream, []byte(`]`)...)

	plansAPI := testServer(planStream)
	plansURL = plansAPI.URL

	fmt.Printf("Applying payment plan:\n%s\n\n", planBytes)

	return plansAPI
}

func applyPayments(amt, freq, n int, sDate time.Time) *httptest.Server {
	var payments [][]byte

	for i := 0; i < n; i++ {
		test := paymentResponse{USD(amt), sDate.Format(ISO8601)}
		testBytes, _ := json.Marshal(test)
		fmt.Printf("Applying payment: %s\n", testBytes)
		payments = append(payments, testBytes)
		sDate.AddDate(0, 0, freq)
	}

	totalPayments := []byte(`[`)
	totalPayments = append(totalPayments, bytes.Join(payments, []byte(`,`))...)
	totalPayments = append(totalPayments, []byte(`]`)...)

	paymentsAPI := testServer(totalPayments)
	paymentsURL = paymentsAPI.URL
	return paymentsAPI
}

func testServer(response []byte) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		resp.Write(response)
	}))
}
