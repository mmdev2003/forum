package status

var CreateStatus = `
INSERT INTO statuses (status_id, account_id, payment_id, payment_status, expiration_at)
VALUES (@status_id, @account_id, @payment_id, @payment_status, @expiration_at)
RETURNING id;
`

var ConfirmPaymentForStatus = `
UPDATE statuses
SET payment_status = @payment_status
WHERE payment_id = @payment_id;
`

var StatusesByAccountID = `
SELECT * FROM statuses
WHERE account_id = @account_id;
`

var DeleteStatus = `
DELETE FROM statuses
WHERE account_id = @account_id and status_id = @status_id;
`
