package message

var CreateMessage = `
INSERT INTO messages (message_owner_account_id, message_owner_login, topic_id, message_reply_to_id, message_text)
VALUES (@message_owner_account_id, @message_owner_login, @topic_id, @message_reply_to_id, @message_text)
RETURNING id;
`

var AddFileToMessage = `
INSERT INTO message_files (message_id, size, url, name, extension)
VALUES (@message_id, @size, @url, @name, @extension)
RETURNING id;
`

var CreateMessageLike = `
INSERT INTO message_likes (topic_id, message_id, like_type_id, account_id)
VALUES (@topic_id, @message_id, @like_type_id, @account_id)
RETURNING id;
`

var IncrementLikeCountToMessage = `
UPDATE messages
SET message_like_count = message_like_count + 1
WHERE id = @message_id;
`

var DecrementLikeCountToMessage = `
UPDATE messages
SET message_like_count = message_like_count - 1
WHERE id = @message_id;
`

var AddReplyCountToMessage = `
UPDATE messages
SET message_reply_count = message_reply_count + 1
WHERE id = @message_id;
`

var AddReportCountToMessage = `
UPDATE messages
SET message_report_count = message_report_count + 1
WHERE id = @message_id;
`

var EditMessage = `
UPDATE messages
SET message_text = @message_text
WHERE id = @message_id;
`

var MessagesByTopicID = `
SELECT * FROM messages
WHERE topic_id = @topic_id;
`

var LikesByTopicIDAndAccountID = `
SELECT * FROM message_likes
WHERE topic_id = @topic_id AND account_id = @account_id;
`

var MessageByAccountID = `
SELECT * FROM messages
WHERE message_owner_account_id = @account_id;
`

var DeleteMessageLike = `
DELETE FROM message_likes
WHERE message_id = @message_id AND account_id = @account_id;
`

var FilesByMessageID = `
SELECT * FROM message_files
WHERE message_id = @message_id;
`
