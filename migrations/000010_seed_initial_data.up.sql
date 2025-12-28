-- Insert default categories
INSERT INTO categories (nama_category, slug, parent_id) VALUES
('Elektronik', 'elektronik', NULL),
('Fashion', 'fashion', NULL),
('Makanan & Minuman', 'makanan-minuman', NULL),
('Kesehatan & Kecantikan', 'kesehatan-kecantikan', NULL),
('Olahraga', 'olahraga', NULL);

-- Insert subcategories
INSERT INTO categories (nama_category, slug, parent_id) VALUES
('Smartphone', 'smartphone', 1),
('Laptop', 'laptop', 1),
('Pakaian Pria', 'pakaian-pria', 2),
('Pakaian Wanita', 'pakaian-wanita', 2),
('Makanan Ringan', 'makanan-ringan', 3);

-- Insert admin user (password should be hashed in real application)
INSERT INTO users (nama, kata_sandi, email, is_admin, status) VALUES
('Admin', 'admin', 'admin@gmail.com', TRUE, 'active');