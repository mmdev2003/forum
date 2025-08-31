package topic

var CreateTopic = `
INSERT INTO topics (thread_id, thread_name, subthread_id, subthread_name, topic_name, topic_owner_account_id, topic_owner_login, topic_moderation_status, topic_is_author)
VALUES (@thread_id, @thread_name, @subthread_id, @subthread_name, @topic_name, @topic_owner_account_id, @topic_owner_login, @topic_moderation_status, @topic_is_author)
RETURNING id;
`

var UpdateTopicLastMessage = `
UPDATE topics
SET topic_last_message_login = @topic_last_message_login, topic_last_message_text = @topic_last_message_text
WHERE id = @topic_id;
`

var RejectTopic = `
UPDATE topics
SET topic_moderation_status = @topic_moderation_status
WHERE id = @topic_id;
`

var ApproveTopic = `
UPDATE topics
SET topic_moderation_status = @topic_moderation_status
WHERE id = @topic_id;
`

var AddMessageCountToTopic = `
UPDATE topics
SET topic_message_count = topic_message_count + 1
WHERE id = @topic_id;
`

var AddViewToTopic = `
UPDATE topics
SET topic_view_count = topic_view_count + 1
WHERE id = @topic_id;
`

var CloseTopic = `
UPDATE topics
SET topic_is_closed = true
WHERE id = @topic_id;
`

var ChangeTopicPriority = `
UPDATE topics
SET topic_priority = @topic_priority
WHERE id = @topic_id;
`

var TopicsBySubthreadID = `
SELECT * FROM topics
WHERE subthread_id = @subthread_id and topic_moderation_status = @topic_moderation_status;
`

var TopicsBySubthreadIDAndAccountID = `
SELECT * FROM topics
WHERE subthread_id = @subthread_id and topic_moderation_status = @topic_moderation_status and topic_owner_account_id = @account_id;
`

var TopicsOnModeration = `
SELECT * FROM topics
WHERE topic_moderation_status = @topic_moderation_status;
`

var TopicsByAccountID = `
SELECT * FROM topics
WHERE topic_owner_account_id = @account_id;
`

var TopicsByAccountIDToday = `
SELECT * FROM topics
WHERE 
    topic_owner_account_id = @account_id
    AND created_at >= NOW() - INTERVAL '24 hours';
`

var UpdateTopicAvatar = `
UPDATE topics
SET topic_avatar_url = @topic_avatar
WHERE id = @topic_id;
`
