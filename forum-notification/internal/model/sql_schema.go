package model

const CreateTableQuery = `
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_type WHERE typname = 'notification_type_enum'
    ) THEN
        CREATE TYPE notification_type_enum AS ENUM (
            'MessageFromTopic',
            'MessageReplyFromTopic',
            'LikeMessageFromTopic',
            'TopicClosed',
            'ResponseToSupportRequest',
            'StatusReceived',
            'FrameReceived',
            'MessageFromDialog',
            'MentionFromTopic',
            'WarningFromAdmin'
        );
    END IF;
END$$;

CREATE TABLE IF NOT EXISTS notification (
    id SERIAL PRIMARY KEY,
    type notification_type_enum NOT NULL,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    account_id INT NOT NULL
);

CREATE TABLE IF NOT EXISTS notification_message_from_topic (
    id SERIAL PRIMARY KEY,
    notification_id INT NOT NULL REFERENCES notification(id) ON DELETE CASCADE,
    message_id INT NOT NULL,
    replier_account_id INT NOT NULL,
    topic_id INT NOT NULL,
    message_text TEXT NOT NULL,
    topic_name VARCHAR(255) NOT NULL,
    replier_login VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS notification_message_reply_from_topic (
    id SERIAL PRIMARY KEY,
    notification_id INT NOT NULL REFERENCES notification(id) ON DELETE CASCADE,
    message_id INT NOT NULL,
    replier_account_id INT NOT NULL,
    topic_id INT NOT NULL,
    message_text TEXT NOT NULL,
    topic_name VARCHAR(255) NOT NULL,
    replier_login VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS notification_like_message_from_topic (
    id SERIAL PRIMARY KEY,
    notification_id INT NOT NULL REFERENCES notification(id) ON DELETE CASCADE,
    message_id INT NOT NULL,
    liker_account_id INT NOT NULL,
    topic_id INT NOT NULL,
    message_text TEXT NOT NULL,
    topic_name VARCHAR(255) NOT NULL,
    liker_login VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS notification_topic_closed (
    id SERIAL PRIMARY KEY,
    notification_id INT NOT NULL REFERENCES notification(id) ON DELETE CASCADE,
    admin_account_id INT NOT NULL,
    topic_id INT NOT NULL,
    topic_name VARCHAR(255) NOT NULL,
    admin_login VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS notification_response_to_support_request (
    id SERIAL PRIMARY KEY,
    notification_id INT NOT NULL REFERENCES notification(id) ON DELETE CASCADE,
    support_request_id INT NOT NULL
);

CREATE TABLE IF NOT EXISTS notification_support_request_closed (
    id SERIAL PRIMARY KEY,
    notification_id INT NOT NULL REFERENCES notification(id) ON DELETE CASCADE,
    support_request_id INT NOT NULL
);

CREATE TABLE IF NOT EXISTS notification_status_received (
    id SERIAL PRIMARY KEY,
    notification_id INT NOT NULL REFERENCES notification(id) ON DELETE CASCADE,
    status_name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS notification_frame_received (
    id SERIAL PRIMARY KEY,
    notification_id INT NOT NULL REFERENCES notification(id) ON DELETE CASCADE,
    frame_name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS notification_message_from_dialog (
    id SERIAL PRIMARY KEY,
    notification_id INT NOT NULL REFERENCES notification(id) ON DELETE CASCADE,
    message_id INT NOT NULL,
    dialog_id INT NOT NULL,
    sender_account_id INT NOT NULL,
    message_text TEXT NOT NULL,
    sender_login VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS notification_mention_from_topic (
    id SERIAL PRIMARY KEY,
    notification_id INT NOT NULL REFERENCES notification(id) ON DELETE CASCADE,
    message_id INT NOT NULL,
    mention_account_id INT NOT NULL,
    message_text TEXT NOT NULL,
    topic_name VARCHAR(255) NOT NULL,
    mention_login VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS notification_warning_from_admin (
    id SERIAL PRIMARY KEY,
    notification_id INT NOT NULL REFERENCES notification(id) ON DELETE CASCADE,
    admin_account_id INT NOT NULL,
    message_text TEXT NOT NULL,
    admin_login VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS notification_settings (
    account_id BIGINT PRIMARY KEY,
    enabled_types notification_type_enum[] NOT NULL DEFAULT '{}'
);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER trigger_update_requests_timestamp BEFORE UPDATE ON notification
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

`

const DropTableQuery = `
DROP TABLE IF EXISTS notification_warning_from_admin;
DROP TABLE IF EXISTS notification_mention_from_topic;
DROP TABLE IF EXISTS notification_message_from_dialog;
DROP TABLE IF EXISTS notification_frame_received;
DROP TABLE IF EXISTS notification_status_received;
DROP TABLE IF EXISTS notification_support_request_closed;
DROP TABLE IF EXISTS notification_response_to_support_request;
DROP TABLE IF EXISTS notification_topic_closed;
DROP TABLE IF EXISTS notification_like_message_from_topic;
DROP TABLE IF EXISTS notification_message_reply_from_topic;
DROP TABLE IF EXISTS notification_message_from_topic;
DROP TABLE IF EXISTS notification;
DROP TABLE IF EXISTS notification_settings;
DROP FUNCTION IF EXISTS update_updated_at_column();
`
