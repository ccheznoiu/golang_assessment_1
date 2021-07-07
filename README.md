 ### Rationale and Approach
This is a demo script intended to be consumed by proficient Go developers. As such the developer intends to follow the strictest conventions while accomplishing the functional requirements (FR's) with a preference for efficiency over vanity. Data validation was handled delicately, see **Assumption 1**. The most expensive operations will be _1 debt API_ call * _n (number of debt objects) payment_plan API_ calls * _n payment API_ calls.

 The "IDE" of choice was the Go Playground. The developer made full use of golang standard package documentation and consulted third-party sources, with a preference for stackoverflow, for more specialized knowledge. The developer was aware of (but has not worked with) the float currency gotcha implicit in this exercise and considered using a third-party package for handling USD calculation, but preferred defining a custom USD type and using standard packages exclusively. Suffice it to say that this turned into quite an expedition and even more time would have gone into researching the nuances of this particular challenge.

 Testing was approached first with the self-assurance of a developer who "thought of everything" and after the first test failed (for performance issues!) the rigidity of a professional tester. More time would certainly have been devoted to writing test cases.

Assimilating the FR's took the most amount of time. FR4 required a few redesigns. There were a surprising number of implications and ambiguities:
##### Assumption 1
The mock data are valid, particularly fractional cents. Data coherence, however, was not taken for granted and was validated for sanity. Data errors were handled light-handedly in the absence of explicit FR's, and in favor of not allowing them to dominate the design. This is not to downplay the seriousness of, e.g. payment `amount`'s of `0`, or unexpected `installment_frequency`'s. The fact that `0` was a valid debt `id` was allowed to factor in, hence this field is an `*int` such that:
```go
// ...
    if debt.ID == nil {
        log.Println("encountered invalid debt object!")
        continue
    }
```
##### Assumption 2
FR 1 makes mention of "all the debt object's fields returned by the API." This was taken to mean "_both_ debt object fields returned by the `debt` API in addition to the new fields from FR's 1, 3 and 4."
##### Assumption 3
`amount` and `amount_to_pay` are **principal** such that a payment plan is known to be complete **only** once all payments found, which represent the entire payment history, have been applied against their applicable principal, i.e. as opposed to when either of these fields <= 0.
##### Assumption 4
Since no mention was made as to whether a payment plan is active _before_ its `start_date`, this is assumed to be false. This is furthermore appealing since "Payments [may be] made on days outside the expected payment schedule," taken to mean _before the_ `start_date`. The payment schedule is assumed to begin at midnight of the `start_date`.
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
"Payments made on days outside the expected payment schedule [...] do not change/delay the payment schedule," only holds significance when granted **Assumption 3**.
