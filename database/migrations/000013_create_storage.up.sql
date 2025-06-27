CREATE TABLE IF NOT EXISTS storage (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    length BIGINT,
    width BIGINT,
    height BIGINT
);