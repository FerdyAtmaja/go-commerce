-- Remove product logs
DELETE FROM log_produk WHERE id_produk IN (1, 3, 5);

-- Remove transaction items
DELETE FROM detail_trx WHERE id_trx IN (1, 2, 3);

-- Remove transactions
DELETE FROM trx WHERE id_user = 2;

-- Remove product photos
DELETE FROM foto_produk WHERE id_produk IN (1, 2, 3, 4, 5, 6);

-- Remove products
DELETE FROM produk WHERE id_toko IN (1, 2);

-- Remove stores
DELETE FROM toko WHERE id_user = 2;

-- Remove addresses
DELETE FROM alamat WHERE id_user = 2;

-- Remove user jambu
DELETE FROM users WHERE email = 'jambu@gmail.com';