CREATE TABLE reward (
    id SERIAL PRIMARY KEY,  -- SERIAL takes care of auto-increment
    reward_address VARCHAR(50) NOT NULL,
    recipient_address VARCHAR(50) NOT NULL,
    transaction_hash VARCHAR(66) NOT NULL,
    token_amount NUMERIC(30, 0) NOT NULL, 
    status SMALLINT NOT NULL DEFAULT 0,    -- 0 for pending, 1 for success, -1 for failed 
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for optimizing queries
CREATE INDEX reward_recipient_address_idx ON reward (recipient_address);
CREATE INDEX reward_transaction_hash_idx ON reward (transaction_hash);
CREATE INDEX reward_status_idx ON reward (status);
CREATE INDEX reward_recipient_address_status_idx ON reward (recipient_address, status);

-- Create a trigger to update 'updated_at' column on update
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_reward_updated_at
BEFORE UPDATE ON reward
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
