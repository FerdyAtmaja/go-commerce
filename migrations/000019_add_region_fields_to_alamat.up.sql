ALTER TABLE alamat 
ADD COLUMN province_id VARCHAR(10) NOT NULL DEFAULT '',
ADD COLUMN city_id VARCHAR(10) NOT NULL DEFAULT '';

-- Add indexes for better performance
CREATE INDEX idx_alamat_province_id ON alamat(province_id);
CREATE INDEX idx_alamat_city_id ON alamat(city_id);