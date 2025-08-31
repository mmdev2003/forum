package model

var CreateTableQuery = `
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
	account_id INTEGER NOT NULL,

	login TEXT NOT NULL,
	header_url TEXT DEFAULT '',
	avatar_url TEXT DEFAULT '',

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users_bans (
    id SERIAL PRIMARY KEY,
    
    from_account_id INTEGER NOT NULL,
    to_account_id INTEGER NOT NULL,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS warnings_from_admins (
    id SERIAL PRIMARY KEY,
    
    admin_account_id INTEGER NOT NULL,
    admin_login TEXT NOT NULL,
    to_account_id INTEGER NOT NULL,
    warning_type TEXT NOT NULL,
    warning_text TEXT NOT NULL,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
   
CREATE TRIGGER update_updated_at_trigger BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER update_updated_at_trigger BEFORE UPDATE ON users_bans
FOR EACH ROW EXECUTE FUNCTION update_updated_at();
`

var DropTableQuery = `
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS users_bans;
`
