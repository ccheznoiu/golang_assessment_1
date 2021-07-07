## Objective
A **Go** script to consume data from REST services serving the objects described below and output debt objects to stdout in JSON Lines format.
##### Debt Object
* `id`
* `amount`
* `is_in_payment_plan`:
    * `true` when the debt is associated with an "active" payment plan
    * `false` when there is no payment plan, or the payment plan is completed
* `remaining_amount`:
    * If the debt is associated with a payment plan, calculated from the payment plan's `amount_to_pay`
* `next_payment_due_date`: 
    * `null` if there is no payment plan or if the debt has been paid off 
    * Calculated by using the payment plan `start_date` and `installment_frequency`. It should be the next installment date after the latest payment, even if this date is in the past
##### Plan Object
* `id`
* `debt_id`
* `amount_to_pay`
* `installment_frequency`: one of: `WEEKLY` or `BI_WEEKLY`
* `installment_amount`
* `start_date`: when the first payment is due
##### Payment Object
* `payment_plan_id`
* `amount`
* `date`
## Rationale
##### A Note on the USD Object
A custom type was preferred for demonstration purposes. `math/big` or a third-party package devoted to currency would have been viable alternatives, except for the observation, discovered in the midst of unit testing, that [$1.xx5 is a special case](https://play.golang.org/p/Hr_WKI1eQLJ) (see [Assumption 1](https://github.com/ccheznoiu/golang_assessment_1/blob/main/README.md#assumption-1)). As such the herein provided `USD` type takes the approach of physically displacing the decimal two places to the right (like my 7th Grade teacher taught me). 
##### Assumption 1
Round thousandth dollar values up (away from zero) to the next hundredth (cent).
##### Assumption 2
Data was validated for sanity. Data errors were handled light-handedly in the absence of explicit FR's, and in favor of not allowing them to dominate the design. This is not to downplay the seriousness of, e.g. payment `amount`'s of `0`, or unexpected `installment_frequency`'s. The fact that `0` was a valid debt `id` was allowed to factor in, hence this field is an `*int` such that:
```go
// ...
    if debt.ID == nil {
        log.Println("encountered invalid debt object!")
        continue
    }
```
##### Assumption 3
`amount` and `amount_to_pay` are **principal** such that a payment plan is known to be complete **only** once all `payment`'s found, which represent the entire payment history, have been applied against their applicable principal, i.e. as opposed to when either of these fields <= 0, i.e. "the calculation."
##### Assumption 4
Since no mention was made as to whether a payment plan is "active" _before_ its `start_date`, this is assumed to be false. This is furthermore appealing since "Payments [may be] made on days outside the expected payment schedule," taken to mean _before the_ `start_date`. The payment schedule is assumed to begin at midnight of the `start_date`.
##### Assumption 5
The payment object field `payment_plan_id` is unfortunately named as it implies that payments may only be made against a payment plan when this is clearly not true. It will be assumed to be 1:1 with the debt `id`.
##### Assumption 6
"the next installment date after the latest payment, even if this date is in the past" seems to mean: "obtain the `next_payment_due_date` by advancing the `start_date` by iterations of the `installment_frequency` while the `next_payment_due_date` is before the latest payment `date`, even if it results in a past date", namely:
```go
func (d *debtResponse) creditPayments(pf int)
// ...
    for nextInstall.Before(pDate) {
        nextInstall = nextInstall.AddDate(0, 0, pf)
    }
```
"Payments made on days outside the expected payment schedule [...] do not change/delay the payment schedule," only holds significance when granted [Assumption 4(https://github.com/ccheznoiu/golang_assessment_1/blob/main/README.md#assumption-4).
