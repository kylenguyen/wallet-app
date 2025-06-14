-- Enable the pgcrypto extension to use gen_random_uuid() for generating UUIDs.
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- =================================================================
--  Define ENUM Types
-- =================================================================

-- Using an ENUM for transaction types provides better data integrity than a string.
CREATE TYPE transaction_type AS ENUM (
    'deposit',
    'withdrawal',
    'transfer'
);


-- =================================================================
--  Create Tables
-- =================================================================

-- Table to store user information.
CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       name VARCHAR(255) NOT NULL,
                       email VARCHAR(255) UNIQUE NOT NULL,
                       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


-- Table to store wallet information for each user.
-- A user can have multiple wallets.
CREATE TABLE wallets (
                         id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                         user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                         name VARCHAR(255) NOT NULL,
                         balance DECIMAL(19, 4) NOT NULL DEFAULT 0.00 CHECK (balance >= 0),
                         created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                         updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


-- Table to store an immutable record of all transactions.
CREATE TABLE transactions (
                              id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                              wallet_id UUID NOT NULL REFERENCES wallets(id),
                              type transaction_type NOT NULL,
                              amount DECIMAL(19, 4) NOT NULL,
                              related_wallet_id UUID NULL REFERENCES wallets(id),
                              created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- =================================================================
--  Create Indexes for Performance
-- =================================================================

-- Index on the user_id in the wallets table for fast lookup of a user's wallets.
CREATE INDEX idx_wallets_user_id ON wallets(user_id);

-- Index on the wallet_id in the transactions table for fast lookup of a wallet's history.
CREATE INDEX idx_transactions_wallet_id ON transactions(wallet_id);

-- Index on the related_wallet_id for efficient querying of transfers.
CREATE INDEX idx_transactions_related_wallet_id ON transactions(related_wallet_id);


-- =================================================================
--  Create Triggers to Automate Timestamps
-- =================================================================

-- A function that updates the updated_at column to the current timestamp.
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;


-- A trigger on the wallets table that calls the function whenever a row is updated.
CREATE TRIGGER set_wallets_updated_at
    BEFORE UPDATE ON wallets
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();