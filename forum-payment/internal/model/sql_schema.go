package model

var CreateTableQuery = `
CREATE TABLE IF NOT EXISTS payments (
    id SERIAL PRIMARY KEY,
    account_id INTEGER NOT NULL,
    product_type TEXT NOT NULL,
	
    address TEXT NOT NULL,
    currency TEXT NOT NULL,
	amount TEXT NOT NULL,
    status TEXT NOT NULL,

    tx_id TEXT DEFAULT '',
    is_paid BOOLEAN DEFAULT FALSE,

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

CREATE TRIGGER update_updated_at_trigger BEFORE UPDATE ON payments
FOR EACH ROW EXECUTE FUNCTION update_updated_at();
`

var DropTableQuery = `
DROP TABLE IF EXISTS payments;
`
