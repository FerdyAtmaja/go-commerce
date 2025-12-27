CREATE TABLE users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    nama VARCHAR(255) NOT NULL,
    kata_sandi VARCHAR(255) NOT NULL,
    notelp VARCHAR(20),
    tanggal_lahir DATE,
    jenis_kelamin ENUM('L', 'P'),
    tentang TEXT,
    pekerjaan VARCHAR(255),
    email VARCHAR(255) UNIQUE NOT NULL,
    id_provinsi BIGINT,
    id_kota BIGINT,
    is_admin BOOLEAN DEFAULT FALSE,
    email_verified_at TIMESTAMP NULL,
    last_login_at TIMESTAMP NULL,
    status ENUM('active', 'blocked') DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_is_admin ON users(is_admin);