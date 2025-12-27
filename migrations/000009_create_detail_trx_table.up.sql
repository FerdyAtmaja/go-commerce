CREATE TABLE detail_trx (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    id_trx BIGINT UNSIGNED NOT NULL,
    id_log_produk BIGINT UNSIGNED NOT NULL,
    id_toko BIGINT UNSIGNED NOT NULL,
    kuantitas INT NOT NULL,
    harga_satuan DECIMAL(12,2) NOT NULL,
    harga_total DECIMAL(14,2) NOT NULL,
    nama_produk_snapshot VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (id_trx) REFERENCES trx(id) ON DELETE CASCADE,
    FOREIGN KEY (id_log_produk) REFERENCES log_produk(id) ON DELETE RESTRICT,
    FOREIGN KEY (id_toko) REFERENCES toko(id) ON DELETE RESTRICT
);

CREATE INDEX idx_detail_trx_trx ON detail_trx(id_trx);
CREATE INDEX idx_detail_trx_log_produk ON detail_trx(id_log_produk);
CREATE INDEX idx_detail_trx_toko ON detail_trx(id_toko);