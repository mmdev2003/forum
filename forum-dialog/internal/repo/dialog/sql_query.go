package dialog

var CreateDialog = `
INSERT INTO dialogs (account1_id, account2_id)
VALUES (@account1_id, @account2_id)
RETURNING id;
`

var CreateFile = `
INSERT INTO message_files (size, url, name, extension)
VALUES (@size, @url, @name, @extension)
RETURNING id;
`

var CreateMessage = `
INSERT INTO messages (dialog_id, from_account_id, to_account_id, message_text)
VALUES (@dialog_id, @from_account_id, @to_account_id, @message_text)
RETURNING id;
`

var AddFileToMessage = `
UPDATE message_files
SET message_id = @message_id
WHERE url = @file_url;
`

var DeleteDialog = `
DELETE FROM dialogs
WHERE id = @dialog_id;
`

var MarkDialogAsStarred = `
UPDATE dialogs
SET 
    is_starred_by_account1 = CASE 
        WHEN account1_id = @account_id THEN NOT is_starred_by_account1 
        ELSE is_starred_by_account1 
    END,
    is_starred_by_account2 = CASE 
        WHEN account2_id = @account_id THEN NOT is_starred_by_account2 
        ELSE is_starred_by_account2 
    END
WHERE id = @dialog_id;
`

var UpdateLastMessageAt = `
UPDATE dialogs
SET last_message_at = @last_message_at
WHERE id = @dialog_id;
`

var MarkMessagesAsRead = `
UPDATE messages
SET is_read = true
WHERE dialog_id = @dialog_id;
`

var DialogsByAccountID = `
SELECT * FROM dialogs 
WHERE account1_id = @account_id OR account2_id = @account_id;
`

var MessagesByDialogId = `
SELECT * FROM messages
WHERE dialog_id = @dialog_id
`

var FilesByMessageID = `
SELECT * FROM message_files
WHERE message_id = @message_id;
`
