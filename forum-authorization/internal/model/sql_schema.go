package model

var CreateTableQuery = `
CREATE TABLE IF NOT EXISTS accounts (
    id SERIAL PRIMARY KEY,
    
    account_id INTEGER NOT NULL UNIQUE,
    refresh_token TEXT DEFAULT '',
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';

CREATE TRIGGER update_updated_at_trigger BEFORE UPDATE ON accounts
FOR EACH ROW EXECUTE FUNCTION update_updated_at();
`

var DropTableQuery = `
DROP TABLE IF EXISTS accounts;
`
