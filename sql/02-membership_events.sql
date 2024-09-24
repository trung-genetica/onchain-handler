CREATE TABLE membership_events (
    id BIGSERIAL PRIMARY KEY,
    user_address VARCHAR(50) NOT NULL,
    order_id BIGINT NOT NULL,
    transaction_hash VARCHAR(66) NOT NULL,
    amount DECIMAL(30, 18),
    status SMALLINT NOT NULL DEFAULT 0,  -- 0 for pending, 1 for success, -1 for failed 
    error_message TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX membership_events_order_id_idx ON membership_events (order_id);

-- Create a trigger to update 'updated_at' column on update
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_membership_events_updated_at
BEFORE UPDATE ON membership_events
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();