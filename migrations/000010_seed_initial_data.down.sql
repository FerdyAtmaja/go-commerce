-- Remove seeded data
DELETE FROM users WHERE email = 'admin@go-commerce.com';
DELETE FROM categories WHERE slug IN ('elektronik', 'fashion', 'makanan-minuman', 'kesehatan-kecantikan', 'olahraga', 'smartphone', 'laptop', 'pakaian-pria', 'pakaian-wanita', 'makanan-ringan');