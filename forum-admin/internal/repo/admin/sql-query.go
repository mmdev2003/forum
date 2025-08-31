package admin

var CreateAdmin = `
INSERT INTO admins (account_id)
VALUES (@account_id)
RETURNING id;
`

var AllAdmin = `
SELECT * FROM admins
`
