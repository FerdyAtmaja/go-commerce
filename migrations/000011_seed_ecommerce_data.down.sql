-- Clean up seeded data in reverse order (due to foreign key constraints)

-- Delete Product Logs
DELETE FROM log_produk WHERE id_produk IN (1,2,3,4,5);

-- Delete Transaction Items
DELETE FROM detail_trx WHERE id_trx IN (1,2,3,4);

-- Delete Transactions
DELETE FROM trx WHERE kode_invoice IN ('INV-2024-001', 'INV-2024-002', 'INV-2024-003', 'INV-2024-004');

-- Delete Product Photos
DELETE FROM foto_produk WHERE id_produk IN (1,2,3,4,5,6,7,8);

-- Delete Products
DELETE FROM produk WHERE slug IN ('iphone-14-pro-max', 'samsung-galaxy-s23-ultra', 'macbook-pro-m2', 'kemeja-formal-pria', 'dress-wanita-elegant', 'serum-vitamin-c', 'protein-whey', 'kopi-arabica-premium');

-- Delete Addresses
DELETE FROM alamat WHERE id_user IN (2,3,4,5);

-- Delete Stores
DELETE FROM toko WHERE nama_toko IN ('Tech Store', 'Fashion Hub', 'Health & Beauty', 'Food Corner');

-- Delete Users (keep admin from 000010)
DELETE FROM users WHERE email IN ('john@example.com', 'jane@example.com', 'bob@example.com', 'alice@example.com');

-- Delete additional Categories (keep original ones from 000010)
DELETE FROM categories WHERE slug IN ('buku-alat-tulis', 'rumah-tangga', 'aksesoris', 'sepatu');