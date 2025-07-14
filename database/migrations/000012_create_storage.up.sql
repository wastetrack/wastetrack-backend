CREATE TABLE IF NOT EXISTS storage (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    length DECIMAL,
    width DECIMAL,
    height DECIMAL,
    is_for_recycled_material BOOLEAN DEFAULT FALSE
);