-- Update trx table to separate payment_status and order_status
ALTER TABLE trx 
ADD COLUMN payment_status ENUM('pending', 'paid', 'failed', 'refunded') DEFAULT 'pending' AFTER metode_bayar,
ADD COLUMN order_status ENUM('created', 'processed', 'shipped', 'delivered', 'cancelled') DEFAULT 'created' AFTER payment_status,
ADD COLUMN shipped_at TIMESTAMP NULL AFTER paid_at;

-- Create indexes for new status columns
CREATE INDEX idx_trx_payment_status ON trx(payment_status);
CREATE INDEX idx_trx_order_status ON trx(order_status);

-- Update existing data to match new schema
UPDATE trx SET 
    payment_status = CASE 
        WHEN status_pembayaran = 'paid' THEN 'paid'
        WHEN status_pembayaran = 'cancelled' THEN 'failed'
        ELSE 'pending'
    END,
    order_status = CASE 
        WHEN status_pembayaran = 'done' THEN 'delivered'
        WHEN status_pembayaran = 'shipped' THEN 'shipped'
        WHEN status_pembayaran = 'cancelled' THEN 'cancelled'
        ELSE 'created'
    END;