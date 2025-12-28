-- Insert user jambu with bcrypt hashed password (password)
INSERT IGNORE INTO users (nama, email, notelp, tanggal_lahir, jenis_kelamin, pekerjaan, is_admin, kata_sandi, created_at, updated_at) VALUES
('Jambu', 'jambu@gmail.com', '081234567890', '1995-07-12', 'L', 'Pengusaha', 0, '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', NOW(), NOW());

-- Insert addresses for jambu (user_id = 6)
INSERT IGNORE INTO alamat (id_user, judul_alamat, nama_penerima, notelp, detail_alamat, kode_pos, is_default, created_at, updated_at) VALUES
(6, 'Rumah Utama', 'Jambu', '081234567890', 'Jl. Merdeka No. 123, RT 01/RW 02, Kelurahan Sukamaju', '12345', 1, NOW(), NOW()),
(6, 'Kantor', 'Jambu', '081234567890', 'Jl. Sudirman No. 456, Lantai 5, Gedung Plaza', '12346', 0, NOW(), NOW()),
(6, 'Rumah Orang Tua', 'Bapak Jambu', '081234567891', 'Jl. Kenangan No. 789, Kampung Halaman', '12347', 0, NOW(), NOW());

-- Insert stores for jambu (user_id = 6)
INSERT IGNORE INTO toko (id_user, nama_toko, deskripsi, url_foto, status, rating, created_at, updated_at) VALUES
(6, 'Jambu Electronics', 'Toko elektronik terpercaya dengan produk berkualitas tinggi', 'https://example.com/jambu-electronics.jpg', 'active', 4.8, NOW(), NOW()),
(6, 'Jambu Fashion Store', 'Fashion trendy untuk anak muda dengan harga terjangkau', 'https://example.com/jambu-fashion.jpg', 'active', 4.6, NOW(), NOW());

-- Insert products for jambu's stores (store_id = 5 & 6)
INSERT IGNORE INTO produk (nama_produk, slug, harga_reseller, harga_konsumen, stok, deskripsi, id_toko, id_category, status, berat, sold_count, created_at, updated_at) VALUES
('iPhone 15 Pro Max', 'iphone-15-pro-max', 18000000, 20000000, 10, 'iPhone terbaru dengan teknologi A17 Pro chip', 5, 6, 'active', 500, 5, NOW(), NOW()),
('MacBook Air M2', 'macbook-air-m2', 15000000, 17000000, 5, 'Laptop tipis dan ringan dengan performa tinggi', 5, 7, 'active', 1200, 3, NOW(), NOW()),
('Samsung Galaxy S24', 'samsung-galaxy-s24', 12000000, 14000000, 15, 'Smartphone Android flagship dengan kamera canggih', 5, 6, 'active', 450, 8, NOW(), NOW()),
('Kemeja Casual Pria', 'kemeja-casual-pria', 150000, 200000, 25, 'Kemeja casual berkualitas untuk pria modern', 6, 8, 'active', 300, 12, NOW(), NOW()),
('Dress Wanita Hebat', 'dress-wanita-hebat', 250000, 350000, 20, 'Dress elegant untuk wanita karir dan acara formal', 6, 9, 'active', 400, 7, NOW(), NOW()),
('Jaket Hoodie Unisex', 'jaket-hoodie-unisex', 180000, 250000, 30, 'Jaket hoodie nyaman untuk pria dan wanita', 6, 8, 'active', 600, 15, NOW(), NOW());

-- Insert product photos
INSERT IGNORE INTO foto_produk (id_produk, url, is_primary, position, created_at, updated_at) VALUES
(9, '/uploads/products/iphone-15-pro-max-1.jpg', 1, 1, NOW(), NOW()),
(9, '/uploads/products/iphone-15-pro-max-2.jpg', 0, 2, NOW(), NOW()),
(10, '/uploads/products/macbook-air-m2-1.jpg', 1, 1, NOW(), NOW()),
(11, '/uploads/products/samsung-galaxy-s24-1.jpg', 1, 1, NOW(), NOW()),
(12, '/uploads/products/kemeja-casual-pria-1.jpg', 1, 1, NOW(), NOW()),
(13, '/uploads/products/dress-wanita-hebat-1.jpg', 1, 1, NOW(), NOW()),
(14, '/uploads/products/jaket-hoodie-unisex-1.jpg', 1, 1, NOW(), NOW());

-- Insert product logs for sold items (must be inserted first)
INSERT IGNORE INTO log_produk (id_produk, nama_produk, slug, harga_konsumen, harga_reseller, deskripsi, id_toko, id_category, created_at, updated_at) VALUES
(9, 'iPhone 15 Pro Max', 'iphone-15-pro-max', 20000000, 18000000, 'iPhone terbaru dengan teknologi A17 Pro chip', 5, 6, NOW(), NOW()),
(13, 'Dress Wanita Hebat', 'dress-wanita-hebat', 350000, 250000, 'Dress elegant untuk wanita karir dan acara formal', 6, 9, NOW(), NOW()),
(11, 'Samsung Galaxy S24', 'samsung-galaxy-s24', 14000000, 12000000, 'Smartphone Android flagship dengan kamera canggih', 5, 6, NOW(), NOW());

-- Insert transactions for jambu (as buyer, user_id = 6)
INSERT IGNORE INTO trx (id_user, alamat_pengiriman, harga_total, kode_invoice, metode_bayar, status_pembayaran, created_at, updated_at) VALUES
(6, 5, 20000000, 'INV-2024-005', 'transfer', 'paid', NOW(), NOW()),
(6, 6, 350000, 'INV-2024-006', 'ewallet', 'paid', NOW(), NOW()),
(6, 7, 14000000, 'INV-2024-007', 'transfer', 'pending', NOW(), NOW());

-- Insert transaction items (using id_log_produk and required fields)
INSERT IGNORE INTO detail_trx (id_trx, id_log_produk, id_toko, kuantitas, harga_satuan, harga_total, nama_produk_snapshot, created_at, updated_at) VALUES
(5, 6, 5, 1, 20000000, 20000000, 'iPhone 15 Pro Max', NOW(), NOW()),
(6, 7, 6, 1, 350000, 350000, 'Dress Wanita Hebat', NOW(), NOW()),
(7, 8, 5, 1, 14000000, 14000000, 'Samsung Galaxy S24', NOW(), NOW());