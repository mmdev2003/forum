package user

var CreateUser = `
INSERT INTO users (account_id, login)
VALUES (@account_id, @login)
RETURNING id;
`

var CreateUserBan = `
INSERT INTO users_bans (from_account_id, to_account_id)
VALUES (@from_account_id, @to_account_id)
RETURNING id;
`

var CreateWarningFromAdmin = `
INSERT INTO warnings_from_admins (admin_account_id, admin_login, to_account_id, warning_type, warning_text)
VALUES (@admin_account_id, @admin_login, @to_account_id, @warning_type, @warning_text)
RETURNING id;
`

var UpdateAvatarUrl = `
UPDATE users
SET avatar_url = @avatar_url
WHERE account_id = @account_id;
`

var UpdateHeaderUrl = `
UPDATE users
SET header_url = @header_url
WHERE account_id = @account_id;
`
var GetUserByAccountID = `
SELECT * FROM users
WHERE account_id = @account_id;
`

var GetUserByLogin = `
SELECT * FROM users
WHERE login = @login;
`

var BanByAccountID = `
SELECT * FROM users_bans
WHERE to_account_id = @to_account_id;
`

var AllWarningFromAdmin = `
SELECT * FROM warnings_from_admins
WHERE to_account_id = @to_account_id;
`

var DeleteUserBan = `
DELETE FROM users_bans
WHERE from_account_id = @from_account_id and to_account_id = @to_account_id;
`
