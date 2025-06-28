Create extension if not exists "uuid-ossp";
CREATE TABLE IF NOT EXISTS storage_items (
    id BIGSERIAL PRIMARY KEY,
    storage_id UUID NOT NULL REFERENCES storage(id) ON DELETE CASCADE,
    waste_type_id UUID NOT NULL REFERENCES waste_types(id) ON DELETE CASCADE,
    weight_kgs DECIMAL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);