-- Add business rule columns to categories table
ALTER TABLE categories ADD COLUMN status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive'));
ALTER TABLE categories ADD COLUMN is_leaf BOOLEAN DEFAULT TRUE;
ALTER TABLE categories ADD COLUMN has_child BOOLEAN DEFAULT FALSE;
ALTER TABLE categories ADD COLUMN has_active_product BOOLEAN DEFAULT FALSE;

-- Add indexes for performance
CREATE INDEX idx_categories_status ON categories(status);
CREATE INDEX idx_categories_is_leaf ON categories(is_leaf);
CREATE INDEX idx_categories_has_child ON categories(has_child);

-- Update existing categories to set proper is_leaf and has_child values
UPDATE categories SET 
    has_child = (SELECT COUNT(*) > 0 FROM categories c2 WHERE c2.parent_id = categories.id),
    is_leaf = (SELECT COUNT(*) = 0 FROM categories c2 WHERE c2.parent_id = categories.id);