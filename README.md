Create a script that will consume data from the HTTP API endpoints and output debts to stdout in JSON Lines format.

Each line should contain:
All the Debt object's fields returned by the API
An additional boolean value, "is_in_payment_plan", which is:
true when the debt is associated with an active payment plan.
false when there is no payment plan, or the payment plan is completed.
Provide a test suite that validates the output being produced, along with any other operations performed internally.
This can be done using any testing technique, but it should provide reasonable coverage of functionality.
Add a new field to the Debt objects in the output: "remaining_amount", containing the calculated amount remaining to be paid on the debt. Output the value as a JSON Number.
If the debt is associated with a payment plan, subtract from the payment plan's amount_to_pay instead. In exchange for signing up for a payment plan, we will allow them to pay a reduced amount to satisfy the debt.
All payments, whether on-time or not, contribute toward paying off a debt.
Add a new field to the Debt object output: "next_payment_due_date", containing the ISO 8601 UTC date (i.e. “2019-09-07”) of when the next payment is due or null if there is no payment plan or if the debt has been paid off.
The next_payment_due_date can be calculated by using the payment plan start_date and installment_frequency. It should be the next installment date after the latest payment, even if this date is in the past.
The next_payment_due_date should be null if there is no payment plan or if the debt has been paid off.
Payments made on days outside the expected payment schedule still go toward paying off the remaining_amount, but do not change/delay the payment schedule.
