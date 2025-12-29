-- Restore payment_status column
ALTER TABLE trx 
ADD COLUMN payment_status ENUM('pending', 'paid', 'failed', 'refunded') DEFAULT 'pending' AFTER metode_bayar;

-- Create index
CREATE INDEX idx_trx_payment_status ON trx(payment_status);

-- Restore data to payment_status from status_pembayaran
UPDATE trx SET 
    payment_status = CASE 
        WHEN status_pembayaran IN ('paid', 'shipped', 'done') THEN 'paid'
        WHEN status_pembayaran = 'failed' THEN 'failed'
        WHEN status_pembayaran = 'refunded' THEN 'refunded'
        ELSE 'pending'
    END;

-- Revert status_pembayaran enum to original
ALTER TABLE trx MODIFY COLUMN status_pembayaran ENUM('pending', 'paid', 'cancelled', 'shipped', 'done') DEFAULT 'pending';