package model

const CreateTableQuery = `
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_type WHERE typname = 'support_request_status'
    ) THEN
        CREATE TYPE support_request_status AS ENUM (
            'open',
            'closed'
        );
    END IF;
END$$;

CREATE TABLE IF NOT EXISTS support_requests (
    id SERIAL PRIMARY KEY,
    account_id INTEGER NOT NULL,
	title VARCHAR(150) NOT NULL,
	description TEXT CHECK (LENGTH(description) <= 5000),
	status support_request_status DEFAULT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS dialogs (
    id SERIAL PRIMARY KEY,
    support_request_id INTEGER NOT NULL REFERENCES support_requests(id),
    user_account_id INTEGER NOT NULL,
    last_message_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT idx_dialogs_support_request 
        UNIQUE (support_request_id, user_account_id)
);

CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    dialog_id INTEGER NOT NULL REFERENCES dialogs(id) ON DELETE CASCADE,
    from_account_id INTEGER NOT NULL,
    to_account_id INTEGER NOT NULL,
    message_text TEXT NOT NULL CHECK (LENGTH(message_text) <= 10000),
    is_read BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT idx_messages_dialog 
        UNIQUE (dialog_id, created_at, id)
);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_requests_timestamp BEFORE UPDATE ON support_requests
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_update_dialog_timestamp BEFORE UPDATE ON dialogs
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_update_message_timestamp BEFORE UPDATE ON messages
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();


CREATE INDEX IF NOT EXISTS idx_support_requests_account_id ON support_requests(account_id);

CREATE INDEX IF NOT EXISTS idx_messages_unread 
ON messages(to_account_id, is_read) 
WHERE is_read = false;


CREATE OR REPLACE FUNCTION update_dialog_last_message()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE dialogs
    SET last_message_at = NEW.created_at
    WHERE id = NEW.dialog_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_dialog_last_message
AFTER INSERT ON messages
FOR EACH ROW EXECUTE FUNCTION update_dialog_last_message();
`

const DropTableQuery = `
DROP TABLE IF EXISTS support_requests CASCADE;
DROP TABLE IF EXISTS dialogs CASCADE;
DROP TABLE IF EXISTS messages CASCADE;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP FUNCTION IF EXISTS update_dialog_last_message();
`
