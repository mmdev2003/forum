package notification

const CreateMessageFromTopicNotification = `
WITH inserted_notification AS (
    INSERT INTO notification (type, account_id)
    VALUES ('MessageFromTopic'::notification_type_enum, @account_id)
    RETURNING id
)
INSERT INTO notification_message_from_topic (notification_id, message_id, replier_account_id, topic_id, message_text, topic_name, replier_login)
SELECT id, @message_id, @replier_account_id, @topic_id, @message_text, @topic_name, @replier_login
FROM inserted_notification
RETURNING id;
`

const CreateMessageReplyFromTopicNotification = `
WITH inserted_notification AS (
    INSERT INTO notification (type, account_id)
    VALUES ('MessageReplyFromTopic'::notification_type_enum, @account_id)
    RETURNING id
)
INSERT INTO notification_message_reply_from_topic (notification_id, message_id, replier_account_id, topic_id, message_text, topic_name, replier_login)
SELECT id, @message_id, @replier_account_id, @topic_id, @message_text, @topic_name, @replier_login
FROM inserted_notification
RETURNING id;
`

const CreateLikeMessageFromTopicNotification = `
WITH inserted_notification AS (
    INSERT INTO notification (type, account_id)
    VALUES ('LikeMessageFromTopic'::notification_type_enum, @account_id)
    RETURNING id
)
INSERT INTO notification_like_message_from_topic (notification_id, message_id, liker_account_id, topic_id, message_text, topic_name, liker_login)
SELECT id, @message_id, @liker_account_id, @topic_id, @message_text, @topic_name, @liker_login
FROM inserted_notification
RETURNING id;
`

const CreateTopicClosedNotification = `
WITH inserted_notification AS (
    INSERT INTO notification (type, account_id)
    VALUES ('TopicClosed'::notification_type_enum, @account_id)
    RETURNING id
)
INSERT INTO notification_topic_closed (notification_id, admin_account_id, topic_id, topic_name, admin_login)
SELECT id, @admin_account_id, @topic_id, @topic_name, @admin_login
FROM inserted_notification
RETURNING id;
`

const CreateResponseToSupportRequestNotification = `
WITH inserted_notification AS (
    INSERT INTO notification (type, account_id)
    VALUES ('ResponseToSupportRequest'::notification_type_enum, @account_id)
    RETURNING id
)
INSERT INTO notification_response_to_support_request (notification_id, support_request_id)
SELECT id, @support_request_id
FROM inserted_notification
RETURNING id;
`

const CreateStatusReceivedNotification = `
WITH inserted_notification AS (
    INSERT INTO notification (type, account_id)
    VALUES ('StatusReceived'::notification_type_enum, @account_id)
    RETURNING id
)
INSERT INTO notification_status_received (notification_id, status_name)
SELECT id, @status_name
FROM inserted_notification
RETURNING id;
`

const CreateFrameReceivedNotification = `
WITH inserted_notification AS (
    INSERT INTO notification (type, account_id)
    VALUES ('FrameReceived'::notification_type_enum, @account_id)
    RETURNING id
)
INSERT INTO notification_frame_received (notification_id, frame_name)
SELECT id, @frame_name
FROM inserted_notification
RETURNING id;
`

const CreateMessageFromDialogNotification = `
WITH inserted_notification AS (
    INSERT INTO notification (type, account_id)
    VALUES ('MessageFromDialog'::notification_type_enum, @account_id)
    RETURNING id
)
INSERT INTO notification_message_from_dialog (notification_id, message_id, dialog_id, sender_account_id, message_text, sender_login)
SELECT id, @message_id, @dialog_id, @sender_account_id, @message_text, @sender_login
FROM inserted_notification
RETURNING id;
`

const CreateMentionFromTopicNotification = `
WITH inserted_notification AS (
    INSERT INTO notification (type, account_id)
    VALUES ('MentionFromTopic'::notification_type_enum, @account_id)
    RETURNING id
)
INSERT INTO notification_mention_from_topic (notification_id, message_id, mention_account_id, message_text, topic_name, mention_login)
SELECT id, @message_id, @mention_account_id, @message_text, @topic_name, @mention_login
FROM inserted_notification
RETURNING id;
`

const CreateWarningFromAdminNotification = `
WITH inserted_notification AS (
    INSERT INTO notification (type, account_id)
    VALUES ('WarningFromAdmin'::notification_type_enum, @account_id)
    RETURNING id
)
INSERT INTO notification_warning_from_admin (notification_id, admin_account_id, message_text, admin_login)
SELECT id, @admin_account_id, @message_text, @admin_login
FROM inserted_notification
RETURNING id;
`

const BaseNotificationSelectQuery = `
SELECT
n.id, n.type, n.is_read, n.created_at, n.updated_at,
n.account_id AS account_id,
COALESCE(mft.message_id, mrf.message_id, lmf.message_id, mfd.message_id, mnt.message_id) AS message_id,
COALESCE(mft.replier_account_id, mrf.replier_account_id) AS replier_account_id,
COALESCE(mft.topic_id, mrf.topic_id, lmf.topic_id, tc.topic_id) AS topic_id,
COALESCE(mft.message_text, mrf.message_text, lmf.message_text, mfd.message_text, mnt.message_text, wa.message_text) AS message_text,
COALESCE(mft.topic_name, mrf.topic_name, lmf.topic_name, tc.topic_name, mnt.topic_name) AS topic_name,
COALESCE(mft.replier_login, mrf.replier_login) AS replier_login,
COALESCE(lmf.liker_account_id, 0) AS liker_account_id,
COALESCE(lmf.liker_login, '') AS liker_login,
COALESCE(tc.admin_account_id, 0) AS admin_account_id,
COALESCE(tc.admin_login, '') AS admin_login,
COALESCE(rsr.support_request_id, 0) AS support_request_id,
COALESCE(src.support_request_id, 0) AS support_request_id_closed,
COALESCE(sr.status_name, '') AS status_name,
COALESCE(fr.frame_name, '') AS frame_name,
COALESCE(mfd.dialog_id, 0) AS dialog_id,
COALESCE(mfd.sender_account_id, 0) AS sender_account_id,
COALESCE(mfd.sender_login, '') AS sender_login,
COALESCE(mnt.mention_account_id, 0) AS mention_account_id,
COALESCE(mnt.mention_login, '') AS mention_login,
COALESCE(wa.admin_account_id, 0) AS warning_admin_account_id,
COALESCE(wa.admin_login, '') AS warning_admin_login
FROM notification n
LEFT JOIN notification_message_from_topic mft ON n.id = mft.notification_id
LEFT JOIN notification_message_reply_from_topic mrf ON n.id = mrf.notification_id
LEFT JOIN notification_like_message_from_topic lmf ON n.id = lmf.notification_id
LEFT JOIN notification_topic_closed tc ON n.id = tc.notification_id
LEFT JOIN notification_response_to_support_request rsr ON n.id = rsr.notification_id
LEFT JOIN notification_support_request_closed src ON n.id = src.notification_id
LEFT JOIN notification_status_received sr ON n.id = sr.notification_id
LEFT JOIN notification_frame_received fr ON n.id = fr.notification_id
LEFT JOIN notification_message_from_dialog mfd ON n.id = mfd.notification_id
LEFT JOIN notification_mention_from_topic mnt ON n.id = mnt.notification_id
LEFT JOIN notification_warning_from_admin wa ON n.id = wa.notification_id
`

const GetNotificationsByAccountID = BaseNotificationSelectQuery + `
WHERE n.account_id = @account_id
ORDER BY n.created_at DESC
`
