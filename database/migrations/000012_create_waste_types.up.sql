CREATE TABLE IF NOT EXISTS waste_types (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    category_id UUID NOT NULL REFERENCES waste_categories(id) ON DELETE CASCADE,
    subcategory_id UUID REFERENCES waste_subcategories(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT
);