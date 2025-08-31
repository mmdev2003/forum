package payment

var CreatePayment = `
INSERT INTO payments (account_id, product_type, currency, amount,  address, status)
VALUES (@account_id, @product_type, @currency, @amount, @address, @status)
RETURNING id;
`

var PaidPayment = `
UPDATE payments
SET is_paid = true
WHERE id = @payment_id;
`

var CancelPayment = `
UPDATE payments
SET status = @status
WHERE id = @payment_id;
`

var ConfirmPayment = `
UPDATE payments
SET status = @status, tx_id = @tx_id
WHERE id = @payment_id;
`

var PaymentByID = `
SELECT * FROM payments
WHERE id = @payment_id;
`

var PaymentByTxID = `
SELECT * FROM payments
WHERE tx_id = @tx_id;
`

var PendingPayments = `
SELECT * FROM payments
WHERE status = @status;
`

var PaymentByAddressAndAmount = `
SELECT * FROM payments
WHERE address = @address AND amount = @amount;
`
