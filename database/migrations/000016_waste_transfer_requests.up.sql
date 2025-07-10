-- Create enum types
DO $$ 
BEGIN
    -- Form types for waste_transfer_forms
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'form_type') THEN
        CREATE TYPE form_type AS ENUM (
            'industry_request', 
            'waste_bank_request'
        );
    END IF;

    -- Status enum for waste_transfer_forms
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'waste_transfer_status') THEN
        CREATE TYPE waste_transfer_status AS ENUM (
            'pending', 
            'assigned', 
            'collecting', 
            'completed', 
            'cancelled',
            'weighing_in',
            'recycling_in_process',
            'weighing_out',
            'recycled',
            'recycle_cancelled'
        );
    END IF;
END $$;

-- Table creation
CREATE TABLE IF NOT EXISTS waste_transfer_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    source_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    destination_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    assigned_collector_id UUID REFERENCES users(id) ON DELETE SET NULL,
    form_type form_type,
    total_weight DECIMAL DEFAULT 0,
    total_price BIGINT DEFAULT 0,
    status waste_transfer_status DEFAULT 'pending',
    image_url TEXT,
    notes TEXT,
    source_phone_number TEXT,
    destination_phone_number TEXT,
    appointment_date DATE,
    appointment_start_time TIMETZ,
    appointment_end_time TIMETZ,
    appointment_location GEOGRAPHY(Point, 4326),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
