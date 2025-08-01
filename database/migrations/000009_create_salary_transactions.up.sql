-- TODO: Revise transaction_type
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create enum types
DO $$ 
BEGIN
    -- Transaction types
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'transaction_type') THEN
        CREATE TYPE transaction_type AS ENUM ('salary', 'waste_payment','point_conversion');
    END IF;
    
    -- Transaction status
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'transaction_status') THEN
        CREATE TYPE transaction_status AS ENUM ('pending','completed', 'failed');
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS salary_transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sender_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    receiver_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    transaction_type transaction_type NOT NULL,
    amount BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    status transaction_status DEFAULT 'completed',
    notes TEXT,
    is_deleted BOOLEAN DEFAULT FALSE
);