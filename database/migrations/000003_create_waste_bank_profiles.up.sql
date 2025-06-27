CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS waste_bank_profiles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    total_waste_weight DECIMAL DEFAULT 0,
    total_workers BIGINT DEFAULT 0,
    open_time TIMETZ,
    close_time TIMETZ
);