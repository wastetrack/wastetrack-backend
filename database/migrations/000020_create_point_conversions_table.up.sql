-- TODO: Revise transaction_type
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create enum types
DO $$ 
BEGIN
    -- Transaction status
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'conversion_status') THEN
        CREATE TYPE conversion_status AS ENUM ('pending','completed', 'rejected');
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS point_conversions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    status conversion_status DEFAULT 'pending',
    is_deleted BOOLEAN DEFAULT FALSE
);