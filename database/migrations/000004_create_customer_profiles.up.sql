CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS customer_profiles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    carbon_deficit BIGINT DEFAULT 0,
    water_saved BIGINT DEFAULT 0,
    bags_stored BIGINT DEFAULT 0,
    trees BIGINT DEFAULT 0
);