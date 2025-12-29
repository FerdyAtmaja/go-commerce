-- Rollback status separation
DROP INDEX idx_trx_payment_status;
DROP INDEX idx_trx_order_status;

ALTER TABLE trx 
DROP COLUMN payment_status,
DROP COLUMN order_status,
DROP COLUMN shipped_at;