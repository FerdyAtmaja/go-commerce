-- Drop indexes first
DROP INDEX IF EXISTS idx_alamat_province_id;
DROP INDEX IF EXISTS idx_alamat_city_id;

-- Drop columns
ALTER TABLE alamat 
DROP COLUMN province_id,
DROP COLUMN city_id;