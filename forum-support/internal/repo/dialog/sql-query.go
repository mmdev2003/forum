package dialog

const CreateDialog = `
INSERT INTO dialogs (support_request_id, user_account_id)
VALUES (@support_request_id, @user_account_id)
RETURNING id;
`

const AddMessageToDialog = `
INSERT INTO messages (dialog_id, from_account_id, to_account_id, message_text)
VALUES (@dialog_id, @from_account_id, @to_account_id, @message_text)
RETURNING id;
`

const MarkMessagesFromRequesterAsRead = `
UPDATE messages
SET is_read = true
WHERE dialog_id = @dialog_id
AND is_read = false
AND from_account_id = @requester_id
`

const MarkMessagesFromSupportAsRead = `
UPDATE messages
SET is_read = true
WHERE dialog_id = @dialog_id
AND is_read = false
AND from_account_id != @requester_id
`

const GetAllDialogs = `
SELECT * FROM dialogs
`

const GetDialogsByAccountID = GetAllDialogs + `
	WHERE user_account_id = @account_id
`

const GetDialogByID = GetAllDialogs + `
	WHERE id = @dialog_id
`

const MessagesByDialogId = `
SELECT * FROM messages
WHERE dialog_id = @dialog_id;
`
