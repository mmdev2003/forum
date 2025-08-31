package model

var CreateTableQuery = `
CREATE TABLE IF NOT EXISTS admins (
    id SERIAL PRIMARY KEY,

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

CREATE TRIGGER update_updated_at_trigger BEFORE UPDATE ON admins
FOR EACH ROW EXECUTE FUNCTION update_updated_at();
`

var DropTableQuery = `
DROP TABLE IF EXISTS admins;
`
