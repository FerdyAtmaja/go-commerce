CREATE TABLE payment_intents (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    trx_id BIGINT UNSIGNED NOT NULL,
    method VARCHAR(50) NOT NULL,
    status ENUM('pending', 'success', 'failed') NOT NULL DEFAULT 'pending',
    expired_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_payment_intents_trx_id (trx_id),
    INDEX idx_payment_intents_status (status),
    INDEX idx_payment_intents_expired_at (expired_at)
);