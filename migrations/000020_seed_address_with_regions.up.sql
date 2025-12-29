-- Update existing addresses with province_id and city_id
-- Jakarta addresses
UPDATE alamat SET province_id = '31', city_id = '3171' WHERE id = 1; -- Jakarta Pusat
UPDATE alamat SET province_id = '31', city_id = '3174' WHERE id = 2; -- Jakarta Selatan  
UPDATE alamat SET province_id = '31', city_id = '3174' WHERE id = 4; -- Jakarta Selatan

-- Bandung address
UPDATE alamat SET province_id = '32', city_id = '3273' WHERE id = 3; -- Bandung

-- Insert additional sample addresses with region data
INSERT IGNORE INTO alamat (id_user, judul_alamat, nama_penerima, notelp, detail_alamat, province_id, city_id, kode_pos, is_default, created_at, updated_at) VALUES
-- Jakarta addresses
(2, 'Kantor', 'John Doe', '081234567891', 'Jl. Thamrin No. 45, Menteng', '31', '3171', '10350', 0, NOW(), NOW()),
(3, 'Rumah Orang Tua', 'Jane Smith', '081234567892', 'Jl. Kemang Raya No. 88, Kemang', '31', '3174', '12560', 0, NOW(), NOW()),

-- Surabaya addresses  
(4, 'Kantor Cabang', 'Bob Wilson', '081234567893', 'Jl. Pemuda No. 123, Gubeng', '35', '3578', '60271', 0, NOW(), NOW()),
(5, 'Rumah Saudara', 'Alice Brown', '081234567894', 'Jl. Diponegoro No. 67, Tegalsari', '35', '3578', '60265', 0, NOW(), NOW()),

-- Yogyakarta addresses
(2, 'Rumah Liburan', 'John Doe', '081234567891', 'Jl. Malioboro No. 15, Gedongtengen', '34', '3471', '55271', 0, NOW(), NOW()),
(3, 'Kos', 'Jane Smith', '081234567892', 'Jl. Kaliurang KM 5, Sleman', '34', '3404', '55281', 0, NOW(), NOW()),

-- Medan addresses
(4, 'Rumah Nenek', 'Bob Wilson', '081234567893', 'Jl. Sisingamangaraja No. 99, Medan Baru', '12', '1275', '20154', 0, NOW(), NOW()),
(5, 'Toko', 'Alice Brown', '081234567894', 'Jl. Kesawan No. 33, Medan Barat', '12', '1271', '20111', 0, NOW(), NOW()),

-- Semarang addresses  
(2, 'Gudang', 'John Doe', '081234567891', 'Jl. Pandanaran No. 77, Semarang Tengah', '33', '3374', '50241', 0, NOW(), NOW()),
(3, 'Rumah Mertua', 'Jane Smith', '081234567892', 'Jl. Gajah Mada No. 44, Semarang Utara', '33', '3374', '50174', 0, NOW(), NOW());