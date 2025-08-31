package authorization

var SetAccount = `
INSERT INTO accounts (account_id)
VALUES (@account_id)
RETURNING id;
`
var UpdateRefreshToken = `
UPDATE accounts 
SET refresh_token = @refresh_token
WHERE account_id = @account_id;
`
var AccountByID = `
SELECT * FROM accounts
WHERE account_id = @account_id
`
var AccountByRefreshToken = `
SELECT * FROM accounts
WHERE refresh_token = @refresh_token;
`
