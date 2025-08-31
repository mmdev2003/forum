package frame

var CreateFrame = `
INSERT INTO frames (frame_id, account_id, payment_id, payment_status, expiration_at)
VALUES (@frame_id, @account_id, @payment_id, @payment_status, @expiration_at)
RETURNING id;
`

var CreateCurrentFrame = `
INSERT INTO current_frame (db_frame_id, account_id)
VALUES (@db_frame_id, @account_id)
RETURNING id;
`

var CreateFrameData = `
INSERT INTO frames_data (monthly_price, forever_price, name, file_id)
VALUES (@monthly_price, @forever_price, @name, @file_id)
RETURNING id;
`

var ConfirmPaymentForFrame = `
UPDATE frames
SET payment_status = @payment_status
WHERE payment_id = @payment_id;
`

var ChangeCurrentFrame = `
UPDATE current_frame
SET db_frame_id = @db_frame_id
WHERE account_id = @account_id;
`

var FramesByAccountID = `
SELECT * FROM frames
WHERE account_id = @account_id;
`

var AllFrames = `
SELECT * FROM frames_data;
`

var GetFrameDataByID = `
SELECT * FROM frames_data
WHERE id = @frame_id;
`

var CurrentFrameByAccountID = `
SELECT * FROM current_frame
WHERE account_id = @account_id;
`
