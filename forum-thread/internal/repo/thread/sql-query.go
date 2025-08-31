package thread

var CreateThread = `
INSERT INTO threads (thread_name, thread_description, allowed_statuses, thread_color)
VALUES (@thread_name, @thread_description, @allowed_statuses, @thread_color)
RETURNING id;
`

var AddMessageCount = `
UPDATE topics
SET topic_message_count = topic_message_count + 1
WHERE id = @topic_id;
`

var AllThreads = `
SELECT * FROM threads;
`
