package model

var CreateTableQuery = `
CREATE TABLE threads (
    id SERIAL PRIMARY KEY,
    
    thread_name TEXT NOT NULL,
    thread_color TEXT NOT NULL,
    thread_description TEXT NOT NULL,
    allowed_statuses TEXT[] NOT NULL,
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE subthreads (
    id SERIAL PRIMARY KEY,
    
    thread_id INTEGER NOT NULL,
    thread_name TEXT NOT NULL,
    
    subthread_name TEXT NOT NULL,
    subthread_description TEXT NOT NULL,
    subthread_message_count INTEGER NOT NULL DEFAULT 0,
    subthread_view_count INTEGER NOT NULL DEFAULT 0,
    subthread_last_message_login TEXT DEFAULT '',
    subthread_last_message_text TEXT DEFAULT '',
    subthread_last_message_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE topics (
    id SERIAL PRIMARY KEY,
    
    thread_id INTEGER NOT NULL,
    thread_name TEXT NOT NULL,
    
    subthread_id INTEGER NOT NULL,
    subthread_name TEXT NOT NULL,
    
    topic_name TEXT NOT NULL,
    topic_owner_account_id INTEGER NOT NULL,
    topic_owner_login TEXT NOT NULL,
    topic_moderation_status TEXT NOT NULL,
    topic_is_author BOOLEAN NOT NULL,
    topic_last_message_login TEXT NOT NULL DEFAULT '',
    topic_last_message_text TEXT NOT NULL DEFAULT '',
    topic_message_count INTEGER NOT NULL DEFAULT 0,
    topic_view_count INTEGER NOT NULL DEFAULT 0,
    topic_is_closed BOOLEAN NOT NULL DEFAULT FALSE,
    topic_priority INTEGER NOT NULL DEFAULT 0,
    topic_avatar_url TEXT NOT NULL DEFAULT '',
    topic_last_message_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    
    topic_id INTEGER NOT NULL,
    
    message_owner_account_id INTEGER NOT NULL,
    message_owner_login TEXT NOT NULL,
    message_text TEXT NOT NULL,
    message_reply_to_id INTEGER NOT NULL,
    message_file_ids INTEGER[] NOT NULL DEFAULT '{}',
    message_like_count INTEGER NOT NULL DEFAULT 0,
    message_reply_count INTEGER NOT NULL DEFAULT 0,
    message_report_count INTEGER NOT NULL DEFAULT 0,
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE message_likes (
    id SERIAL PRIMARY KEY,
    
    topic_id INTEGER NOT NULL,
    
    message_id INTEGER NOT NULL,
    account_id INTEGER NOT NULL,
    like_type_id INTEGER NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE message_files (
    id SERIAL PRIMARY KEY,
    
    message_id INTEGER NOT NULL,
	url TEXT NOT NULL,
	name TEXT NOT NULL,
	size INTEGER NOT NULL,
	extension TEXT NOT NULL,
	
	created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE account_statistics (
    id SERIAL PRIMARY KEY,
    account_id INTEGER NOT NULL UNIQUE,
    created_topics_count INTEGER NOT NULL DEFAULT 0,
    sent_messages_to_topics_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_threads_updated_at BEFORE UPDATE ON threads
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_topics_updated_at BEFORE UPDATE ON topics
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_subthreads_updated_at BEFORE UPDATE ON subthreads
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_messages_updated_at BEFORE UPDATE ON messages
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_message_likes_updated_at BEFORE UPDATE ON message_likes
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_message_files_updated_at BEFORE UPDATE ON message_files
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_account_statistics_updated_at BEFORE UPDATE ON account_statistics
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE INDEX idx_subthreads_thread_id ON subthreads (thread_id);
CREATE INDEX idx_topics_subthread_id ON topics (subthread_id);
CREATE INDEX idx_messages_topic_id ON messages (topic_id);
CREATE INDEX idx_message_likes_message_id ON message_likes (message_id);
CREATE INDEX idx_account_statistics_account_id ON account_statistics (account_id);
CREATE INDEX idx_message_files_account_id ON message_files (message_id);
`

var DropTableQuery = `
DROP TRIGGER IF EXISTS update_threads_updated_at ON threads;
DROP TRIGGER IF EXISTS update_subthreads_updated_at ON subthreads;
DROP TRIGGER IF EXISTS update_topics_updated_at ON topics;
DROP TRIGGER IF EXISTS update_messages_updated_at ON messages;
DROP TRIGGER IF EXISTS update_message_likes_updated_at ON message_likes;
DROP TRIGGER IF EXISTS update_account_statistics_updated_at ON account_statistics;
DROP TRIGGER IF EXISTS update_message_files_updated_at ON message_files;

DROP INDEX IF EXISTS idx_subthreads_thread_id;
DROP INDEX IF EXISTS idx_topics_subthread_id;
DROP INDEX IF EXISTS idx_messages_topic_id;
DROP INDEX IF EXISTS idx_message_likes_message_id;
DROP INDEX IF EXISTS idx_account_statistics_account_id;
DROP INDEX IF EXISTS idx_message_files_account_id;

DROP TABLE IF EXISTS threads;
DROP TABLE IF EXISTS subthreads;
DROP TABLE IF EXISTS topics;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS message_likes;
DROP TABLE IF EXISTS account_statistics;
DROP TABLE IF EXISTS message_files;
`
