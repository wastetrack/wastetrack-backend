CREATE TABLE IF NOT EXISTS waste_bank_priced_types (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    waste_bank_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    waste_type_id UUID NOT NULL REFERENCES waste_types(id) ON DELETE CASCADE,
    custom_price_per_kgs BIGINT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(waste_bank_id, waste_type_id)
);