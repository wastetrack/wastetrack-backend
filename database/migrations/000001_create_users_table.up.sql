CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS postgis;
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
    username TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    role user_role NOT NULL DEFAULT 'customer',
    avatar_url TEXT,
    phone_number TEXT,
    institution TEXT,
    address TEXT,
    city TEXT,
    province TEXT,
    points BIGINT DEFAULT 0,
    balance BIGINT DEFAULT 0,
    location GEOGRAPHY(Point, 4326),
    is_email_verified BOOLEAN DEFAULT FALSE,
    is_accepting_customer BOOLEAN DEFAULT TRUE,
    is_agreed_to_terms BOOLEAN DEFAULT TRUE,
    email_verification_token TEXT,
    email_change_token TEXT,
    email_change_expiry TIMESTAMPTZ,
    new_email TEXT,
    reset_password_token TEXT,
    reset_password_expiry TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);