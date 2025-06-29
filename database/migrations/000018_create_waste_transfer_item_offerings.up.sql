CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS waste_transfer_item_offerings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    transfer_form_id UUID NOT NULL REFERENCES waste_transfer_forms(id) ON DELETE CASCADE,
    waste_type_id UUID NOT NULL REFERENCES waste_types(id) ON DELETE CASCADE,
    offering_weight DECIMAL,
    offering_price_per_kgs BIGINT,
    accepted_weight DECIMAL,
    accepted_price_per_kgs BIGINT
);