-- Seed additional categories (avoiding duplicates from 000010)
INSERT IGNORE INTO categories (nama_category, slug, parent_id, created_at, updated_at) VALUES
('Buku & Alat Tulis', 'buku-alat-tulis', NULL, NOW(), NOW()),
('Rumah Tangga', 'rumah-tangga', NULL, NOW(), NOW()),
('Aksesoris', 'aksesoris', 1, NOW(), NOW()),
('Sepatu', 'sepatu', 2, NOW(), NOW());

-- Seed users (avoiding duplicate admin from 000010)
INSERT IGNORE INTO users (nama, email, notelp, tanggal_lahir, jenis_kelamin, pekerjaan, is_admin, kata_sandi, created_at, updated_at) VALUES
('John Doe', 'john@example.com', '081234567891', '1985-05-15', 'L', 'Software Engineer', 0, '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', NOW(), NOW()),
('Jane Smith', 'jane@example.com', '081234567892', '1990-08-20', 'P', 'Designer', 0, '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', NOW(), NOW()),
('Bob Wilson', 'bob@example.com', '081234567893', '1988-12-10', 'L', 'Marketing', 0, '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', NOW(), NOW()),
('Alice Brown', 'alice@example.com', '081234567894', '1992-03-25', 'P', 'Teacher', 0, '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', NOW(), NOW());

-- Seed stores (using user IDs from seeded users)
INSERT IGNORE INTO toko (id_user, nama_toko, deskripsi, url_foto, status, rating, created_at, updated_at) VALUES
(2, 'Tech Store', 'Toko elektronik terpercaya dengan produk berkualitas', 'https://example.com/tech-store.jpg', 'active', 4.5, NOW(), NOW()),
(3, 'Fashion Hub', 'Koleksi fashion terlengkap untuk pria dan wanita', 'https://example.com/fashion-hub.jpg', 'active', 4.2, NOW(), NOW()),
(4, 'Health & Beauty', 'Produk kesehatan dan kecantikan original', 'https://example.com/health-beauty.jpg', 'active', 4.7, NOW(), NOW()),
(5, 'Food Corner', 'Makanan dan minuman segar setiap hari', 'https://example.com/food-corner.jpg', 'active', 4.3, NOW(), NOW());

-- Seed addresses
INSERT IGNORE INTO alamat (id_user, judul_alamat, nama_penerima, notelp, detail_alamat, kode_pos, is_default, created_at, updated_at) VALUES
(2, 'Rumah', 'John Doe', '081234567891', 'Jl. Merdeka No. 123, Menteng, Jakarta Pusat, DKI Jakarta', '10310', 1, NOW(), NOW()),
(3, 'Kantor', 'Jane Smith', '081234567892', 'Jl. Sudirman No. 456, Kuningan, Jakarta Selatan, DKI Jakarta', '12920', 1, NOW(), NOW()),
(4, 'Rumah', 'Bob Wilson', '081234567893', 'Jl. Gatot Subroto No. 789, Dago, Bandung, Jawa Barat', '40135', 1, NOW(), NOW()),
(5, 'Apartemen', 'Alice Brown', '081234567894', 'Jl. HR Rasuna Said No. 321, Setiabudi, Jakarta Selatan, DKI Jakarta', '12940', 1, NOW(), NOW());

-- Seed products (using existing category IDs and store IDs)
INSERT IGNORE INTO produk (nama_produk, slug, harga_reseller, harga_konsumen, stok, deskripsi, id_toko, id_category, status, berat, sold_count, created_at, updated_at) VALUES
('iPhone 14 Pro Max', 'iphone-14-pro-max', 15000000, 17000000, 50, 'iPhone terbaru dengan teknologi A16 Bionic chip', 1, 6, 'active', 240, 25, NOW(), NOW()),
('Samsung Galaxy S23 Ultra', 'samsung-galaxy-s23-ultra', 14000000, 16000000, 30, 'Smartphone flagship Samsung dengan S Pen', 1, 6, 'active', 234, 18, NOW(), NOW()),
('MacBook Pro M2', 'macbook-pro-m2', 25000000, 28000000, 20, 'Laptop profesional dengan chip M2 terbaru', 1, 7, 'active', 1600, 12, NOW(), NOW()),
('Kemeja Formal Pria', 'kemeja-formal-pria', 150000, 200000, 100, 'Kemeja formal berkualitas untuk pria', 2, 8, 'active', 300, 45, NOW(), NOW()),
('Dress Wanita Elegant', 'dress-wanita-elegant', 250000, 350000, 75, 'Dress elegant untuk acara formal', 2, 9, 'active', 400, 32, NOW(), NOW()),
('Serum Vitamin C', 'serum-vitamin-c', 80000, 120000, 200, 'Serum wajah dengan vitamin C untuk kulit cerah', 3, 4, 'active', 50, 67, NOW(), NOW()),
('Protein Whey', 'protein-whey', 400000, 500000, 80, 'Suplemen protein untuk fitness dan gym', 3, 5, 'active', 1000, 28, NOW(), NOW()),
('Kopi Arabica Premium', 'kopi-arabica-premium', 75000, 100000, 150, 'Kopi arabica pilihan dengan cita rasa premium', 4, 3, 'active', 250, 89, NOW(), NOW());

-- Seed product photos
INSERT IGNORE INTO foto_produk (id_produk, url, is_primary, position, created_at, updated_at) VALUES
(1, '/uploads/products/iphone-14-pro-max-1.jpg', 1, 1, NOW(), NOW()),
(1, '/uploads/products/iphone-14-pro-max-2.jpg', 0, 2, NOW(), NOW()),
(2, '/uploads/products/samsung-s23-ultra-1.jpg', 1, 1, NOW(), NOW()),
(2, '/uploads/products/samsung-s23-ultra-2.jpg', 0, 2, NOW(), NOW()),
(3, '/uploads/products/macbook-pro-m2-1.jpg', 1, 1, NOW(), NOW()),
(4, '/uploads/products/kemeja-formal-1.jpg', 1, 1, NOW(), NOW()),
(5, '/uploads/products/dress-elegant-1.jpg', 1, 1, NOW(), NOW()),
(6, '/uploads/products/serum-vitamin-c-1.jpg', 1, 1, NOW(), NOW()),
(7, '/uploads/products/protein-whey-1.jpg', 1, 1, NOW(), NOW()),
(8, '/uploads/products/kopi-arabica-1.jpg', 1, 1, NOW(), NOW());

-- Seed product logs (must be created first for detail_trx reference)
INSERT IGNORE INTO log_produk (id_produk, nama_produk, slug, harga_konsumen, harga_reseller, deskripsi, id_toko, id_category, created_at, updated_at) VALUES
(1, 'iPhone 14 Pro Max', 'iphone-14-pro-max', 17000000, 15000000, 'iPhone terbaru dengan teknologi A16 Bionic chip', 1, 6, NOW(), NOW()),
(2, 'Samsung Galaxy S23 Ultra', 'samsung-galaxy-s23-ultra', 16000000, 14000000, 'Smartphone flagship Samsung dengan S Pen', 1, 6, NOW(), NOW()),
(3, 'MacBook Pro M2', 'macbook-pro-m2', 28000000, 25000000, 'Laptop profesional dengan chip M2 terbaru', 1, 7, NOW(), NOW()),
(4, 'Kemeja Formal Pria', 'kemeja-formal-pria', 200000, 150000, 'Kemeja formal berkualitas untuk pria', 2, 8, NOW(), NOW()),
(5, 'Dress Wanita Elegant', 'dress-wanita-elegant', 350000, 250000, 'Dress elegant untuk acara formal', 2, 9, NOW(), NOW());

-- Seed transactions
INSERT IGNORE INTO trx (id_user, alamat_pengiriman, harga_total, kode_invoice, metode_bayar, status_pembayaran, created_at, updated_at) VALUES
(2, 1, 17000000, 'INV-2024-001', 'transfer', 'paid', NOW(), NOW()),
(3, 2, 350000, 'INV-2024-002', 'ewallet', 'paid', NOW(), NOW()),
(4, 3, 620000, 'INV-2024-003', 'transfer', 'pending', NOW(), NOW()),
(5, 4, 200000, 'INV-2024-004', 'cod', 'pending', NOW(), NOW());

-- Seed transaction items (using id_log_produk and adding required fields)
INSERT IGNORE INTO detail_trx (id_trx, id_log_produk, id_toko, kuantitas, harga_satuan, harga_total, nama_produk_snapshot, created_at, updated_at) VALUES
(1, 1, 1, 1, 17000000, 17000000, 'iPhone 14 Pro Max', NOW(), NOW()),
(2, 5, 2, 1, 350000, 350000, 'Dress Wanita Elegant', NOW(), NOW()),
(3, 2, 1, 1, 16000000, 16000000, 'Samsung Galaxy S23 Ultra', NOW(), NOW()),
(3, 3, 1, 1, 28000000, 28000000, 'MacBook Pro M2', NOW(), NOW()),
(4, 4, 2, 1, 200000, 200000, 'Kemeja Formal Pria', NOW(), NOW());