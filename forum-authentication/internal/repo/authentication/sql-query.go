package authentication

var CreateAccount = `
INSERT INTO accounts (login, email, password, role)
VALUES (@login, @email, @password, @role)
RETURNING id;
`

var AccountByLogin = `
SELECT * FROM accounts
WHERE login = @login;
`
var AccountByID = `
SELECT * FROM accounts
WHERE id = @account_id;
`

var SetTwoFaKey = `
UPDATE accounts
SET two_fa_key = @two_fa_key
WHERE id = @account_id;
`

var UpgradeToAdmin = `
UPDATE accounts
SET role = @role
WHERE id = @account_id;
`

var UpgradeToSupport = `
UPDATE accounts
SET role = @role
WHERE id = @account_id;
`

var UpdatePassword = `
UPDATE accounts
SET password = @new_password, last_change_password_at = @last_change_password_at
WHERE id = @account_id;
`

var DeleteTwoFaKey = `
UPDATE accounts
SET two_fa_key = ''
WHERE id = @account_id;
`
