CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS waste_drop_request_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    request_id UUID NOT NULL REFERENCES waste_drop_requests(id) ON DELETE CASCADE,
    waste_type_id UUID NOT NULL REFERENCES waste_types(id) ON DELETE CASCADE,
    quantity BIGINT,
    verified_weight DECIMAL,
    verified_price_per_kgs BIGINT,
    verified_subtotal BIGINT,
    is_deleted BOOLEAN DEFAULT FALSE
);