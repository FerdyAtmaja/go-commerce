-- Remove payment_status column and update status_pembayaran enum
-- First, update status_pembayaran to include all payment statuses
ALTER TABLE trx MODIFY COLUMN status_pembayaran ENUM('pending', 'paid', 'failed', 'refunded', 'cancelled', 'shipped', 'done') DEFAULT 'pending';

-- Migrate data from payment_status to status_pembayaran if needed
UPDATE trx SET 
    status_pembayaran = CASE 
        WHEN payment_status = 'paid' AND order_status = 'delivered' THEN 'done'
        WHEN payment_status = 'paid' AND order_status = 'shipped' THEN 'shipped'
        WHEN payment_status = 'paid' THEN 'paid'
        WHEN payment_status = 'failed' THEN 'failed'
        WHEN payment_status = 'refunded' THEN 'refunded'
        ELSE 'pending'
    END
WHERE payment_status IS NOT NULL;

-- Drop payment_status column and its index
DROP INDEX idx_trx_payment_status ON trx;
ALTER TABLE trx DROP COLUMN payment_status;