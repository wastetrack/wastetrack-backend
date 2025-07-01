-- TODO: Revise Request status
DO $$ 
BEGIN
    -- Delivery types
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'delivery_type') THEN
        CREATE TYPE delivery_type AS ENUM ('pickup', 'dropoff');
    END IF;
    
    -- Request status
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'request_status') THEN
        CREATE TYPE request_status AS ENUM ('pending', 'scheduled', 'completed', 'cancelled');
    END IF;
END $$;


CREATE TABLE IF NOT EXISTS waste_drop_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    delivery_type delivery_type NOT NULL,
    customer_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_phone_number TEXT,
    waste_bank_id UUID REFERENCES users(id) ON DELETE SET NULL,
    assigned_collector_id UUID REFERENCES users(id) ON DELETE SET NULL,
    total_price BIGINT DEFAULT 0,
    status request_status DEFAULT 'pending',
    appointment_date DATE,
    appointment_start_time TIMETZ,
    appointment_end_time TIMETZ,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);