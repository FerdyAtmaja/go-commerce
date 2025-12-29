-- Remove additional seeded addresses (keep original 4 addresses)
DELETE FROM alamat WHERE id > 4;

-- Reset province_id and city_id for original addresses
UPDATE alamat SET province_id = '', city_id = '' WHERE id <= 4;