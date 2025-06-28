CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- Create enum types
DO $$ 
BEGIN
    -- Collector management status
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'collector_status') THEN
        CREATE TYPE collector_status AS ENUM ('active', 'inactive');
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS collector_managements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    waste_bank_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    collector_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status collector_status DEFAULT 'active',
    UNIQUE(waste_bank_id, collector_id)
);