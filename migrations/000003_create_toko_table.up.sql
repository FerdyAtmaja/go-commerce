CREATE TABLE toko (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY UNIQUE,
    id_user BIGINT UNSIGNED NOT NULL UNIQUE,
    nama_toko VARCHAR(255) NOT NULL,
    url_foto VARCHAR(255),
    deskripsi TEXT,
    status ENUM('active', 'suspended') DEFAULT 'active',
    rating DECIMAL(2,1) DEFAULT 0.0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (id_user) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_toko_user ON toko(id_user);
CREATE INDEX idx_toko_status ON toko(status);
CREATE INDEX idx_toko_rating ON toko(rating);
CREATE INDEX idx_toko_deleted_at ON toko(deleted_at);