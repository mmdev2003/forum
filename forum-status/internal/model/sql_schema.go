package model

var CreateTableQuery = `
CREATE TABLE IF NOT EXISTS statuses (
    id SERIAL PRIMARY KEY,
    
    status_id INTEGER NOT NULL,
    account_id INTEGER NOT NULL,
    payment_id INTEGER NOT NULL,
    payment_status TEXT NOT NULL,
   
    expiration_at TIMESTAMPTZ NOT NULL,
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

CREATE TRIGGER update_updated_at_trigger BEFORE UPDATE ON statuses
FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE INDEX idx_statuses_account_id ON statuses (account_id);
CREATE INDEX idx_statuses_payment_id ON statuses (payment_id);
`

var DropTableQuery = `
DROP TABLE IF EXISTS statuses;
`
