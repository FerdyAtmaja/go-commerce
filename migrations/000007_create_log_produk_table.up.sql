CREATE TABLE log_produk (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    id_produk BIGINT UNSIGNED NOT NULL,
    nama_produk VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    harga_reseller DECIMAL(12,2) NOT NULL,
    harga_konsumen DECIMAL(12,2) NOT NULL,
    deskripsi TEXT,
    id_toko BIGINT NOT NULL,
    id_category BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_log_produk_produk ON log_produk(id_produk);
CREATE INDEX idx_log_produk_toko ON log_produk(id_toko);
CREATE INDEX idx_log_produk_created ON log_produk(created_at);