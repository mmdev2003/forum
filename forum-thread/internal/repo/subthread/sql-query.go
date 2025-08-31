package subthread

var CreateSubthread = `
INSERT INTO subthreads (thread_id, thread_name, subthread_name, subthread_description)
VALUES (@thread_id, @thread_name, @subthread_name, @subthread_description)
RETURNING id;
`

var AddViewToSubthread = `
UPDATE subthreads
SET subthread_view_count = subthread_view_count + 1
WHERE id = @subthread_id;
`

var AddMessageCountToSubthread = `
UPDATE subthreads
SET subthread_message_count = subthread_message_count + 1
WHERE id = @subthread_id;
`

var UpdateSubthreadLastMessage = `
UPDATE subthreads
SET subthread_last_message_login = @subthread_last_message_login, subthread_last_message_text = @subthread_last_message_text
WHERE id = @subthread_id;
`

var SubthreadsByThreadID = `
SELECT * FROM subthreads
WHERE thread_id = @thread_id;
`
