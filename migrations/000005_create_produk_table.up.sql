CREATE TABLE produk (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    nama_produk VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    harga_reseller DECIMAL(12,2) NOT NULL,
    harga_konsumen DECIMAL(12,2) NOT NULL,
    stok INT DEFAULT 0,
    deskripsi TEXT,
    id_toko BIGINT UNSIGNED NOT NULL,
    id_category BIGINT UNSIGNED NOT NULL,
    status ENUM('active', 'inactive') DEFAULT 'active',
    berat INT DEFAULT 0,
    sold_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (id_toko) REFERENCES toko(id) ON DELETE CASCADE,
    FOREIGN KEY (id_category) REFERENCES categories(id) ON DELETE RESTRICT
);

CREATE INDEX idx_produk_toko ON produk(id_toko);
CREATE INDEX idx_produk_category ON produk(id_category);
CREATE INDEX idx_produk_status ON produk(status);
CREATE INDEX idx_produk_slug ON produk(slug);
CREATE INDEX idx_produk_harga_konsumen ON produk(harga_konsumen);
CREATE INDEX idx_produk_deleted_at ON produk(deleted_at);