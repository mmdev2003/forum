package model

var CreateTableQuery = `
CREATE TABLE IF NOT EXISTS dialogs (
    id SERIAL PRIMARY KEY,
    
    account1_id INTEGER NOT NULL,
    account2_id INTEGER NOT NULL,
	is_starred_by_account1 BOOLEAN DEFAULT FALSE,
	is_starred_by_account2 BOOLEAN DEFAULT FALSE,
    
    last_message_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    dialog_id INTEGER NOT NULL REFERENCES dialogs(id),
    
    from_account_id INTEGER NOT NULL,
    to_account_id INTEGER NOT NULL,
    message_text TEXT NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    file_id INTEGER DEFAULT 0,
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS message_files (
    id SERIAL PRIMARY KEY,
    
    message_id INTEGER DEFAULT 0,
	url TEXT NOT NULL,
	name TEXT NOT NULL,
	size INTEGER NOT NULL,
	extension TEXT NOT NULL,
	
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

CREATE TRIGGER trigger_update_dialogs_timestamp BEFORE UPDATE ON dialogs
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_update_messages_timestamp BEFORE UPDATE ON messages
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_message_files_updated_at BEFORE UPDATE ON message_files
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE INDEX idx_dialogs_account1_id ON dialogs(account1_id);
CREATE INDEX idx_dialogs_account2_id ON dialogs(account2_id);
CREATE INDEX idx_messages_dialog_id ON messages(dialog_id);
CREATE INDEX idx_message_files_account_id ON message_files (message_id);
`

var DropTableQuery = `
DROP TABLE IF EXISTS messages CASCADE;
DROP TABLE IF EXISTS dialogs CASCADE;
DROP TABLE IF EXISTS message_files CASCADE;

DROP FUNCTION IF EXISTS update_updated_at_column();
`
