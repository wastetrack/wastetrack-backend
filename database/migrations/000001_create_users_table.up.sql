CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- Create enum types
DO $$ 
BEGIN
    -- User roles
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role') THEN
        CREATE TYPE user_role AS ENUM ('admin', 'waste_bank_unit', 'waste_collector_unit','waste_bank_central','waste_collector_central', 'customer', 'industry', 'government');
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    role user_role NOT NULL,
    is_email_verified BOOLEAN DEFAULT FALSE,
    phone_number TEXT,
    institution TEXT,
    address TEXT,
    city TEXT,
    province TEXT,
    points DECIMAL DEFAULT 0,
    balance BIGINT DEFAULT 0,
    token TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);