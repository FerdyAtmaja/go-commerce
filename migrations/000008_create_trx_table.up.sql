CREATE TABLE trx (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    id_user BIGINT UNSIGNED NOT NULL,
    alamat_pengiriman BIGINT UNSIGNED NOT NULL,
    harga_total DECIMAL(14,2) NOT NULL,
    kode_invoice VARCHAR(255) UNIQUE NOT NULL,
    metode_bayar ENUM('transfer', 'cod', 'ewallet', 'credit_card'),
    status_pembayaran ENUM('pending', 'paid', 'cancelled', 'shipped', 'done') DEFAULT 'pending',
    paid_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (id_user) REFERENCES users(id) ON DELETE RESTRICT,
    FOREIGN KEY (alamat_pengiriman) REFERENCES alamat(id) ON DELETE RESTRICT
);

CREATE INDEX idx_trx_user ON trx(id_user);
CREATE INDEX idx_trx_status_pembayaran ON trx(status_pembayaran);
CREATE INDEX idx_trx_invoice ON trx(kode_invoice);
CREATE INDEX idx_trx_created ON trx(created_at);