-- Remove business rule columns from categories table
DROP INDEX IF EXISTS idx_categories_has_child ON categories;
DROP INDEX IF EXISTS idx_categories_is_leaf ON categories;
DROP INDEX IF EXISTS idx_categories_status ON categories;

ALTER TABLE categories DROP COLUMN IF EXISTS has_active_product;
ALTER TABLE categories DROP COLUMN IF EXISTS has_child;
ALTER TABLE categories DROP COLUMN IF EXISTS is_leaf;
ALTER TABLE categories DROP COLUMN IF EXISTS status;

-- Remove status column from products table
ALTER TABLE produk DROP COLUMN IF EXISTS status;