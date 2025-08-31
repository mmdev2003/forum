package model

var CreateTableQuery = `
CREATE TABLE IF NOT EXISTS frames (
    id SERIAL PRIMARY KEY,
    
    frame_id INTEGER NOT NULL,
    account_id INTEGER NOT NULL,
    payment_id INTEGER NOT NULL,
    payment_status TEXT NOT NULL,
   
    expiration_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS frames_data (
    id SERIAL PRIMARY KEY,
    
    name TEXT NOT NULL,
	monthly_price FLOAT NOT NULL,
	forever_price FLOAT NOT NULL,
	file_id TEXT NOT NULL,
	
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS current_frame (
    id SERIAL PRIMARY KEY,
    
    db_frame_id INTEGER NOT NULL,
    account_id INTEGER NOT NULL,
    
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

CREATE TRIGGER update_updated_at_trigger BEFORE UPDATE ON frames
FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE INDEX idx_statuses_account_id ON frames (account_id);
CREATE INDEX idx_statuses_payment_id ON frames (payment_id);
`

var DropTableQuery = `
DROP TABLE IF EXISTS frames;
DROP TABLE IF EXISTS frames_data;
DROP TABLE IF EXISTS current_frame;
`
