CREATE TABLE foto_produk (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    id_produk BIGINT UNSIGNED NOT NULL,
    url VARCHAR(255) NOT NULL,
    is_primary BOOLEAN DEFAULT FALSE,
    position BIGINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (id_produk) REFERENCES produk(id) ON DELETE CASCADE
);

CREATE INDEX idx_foto_produk_produk ON foto_produk(id_produk);
CREATE INDEX idx_foto_produk_primary ON foto_produk(is_primary);
CREATE INDEX idx_foto_produk_position ON foto_produk(position);