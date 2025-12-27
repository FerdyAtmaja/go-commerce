CREATE TABLE alamat (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    id_user BIGINT UNSIGNED NOT NULL,
    judul_alamat VARCHAR(255) NOT NULL,
    nama_penerima VARCHAR(255) NOT NULL,
    notelp VARCHAR(20),
    detail_alamat TEXT NOT NULL,
    kode_pos VARCHAR(10),
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (id_user) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_alamat_user ON alamat(id_user);
CREATE INDEX idx_alamat_default ON alamat(is_default);