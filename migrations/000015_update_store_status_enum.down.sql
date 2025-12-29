-- Revert store status enum back to original values
-- First update any new status values to valid old ones
UPDATE toko SET status = 'active' WHERE status IN ('pending', 'inactive');

-- Revert enum to original values
ALTER TABLE toko MODIFY COLUMN status ENUM('active', 'suspended') DEFAULT 'active';

-- Drop the added index
DROP INDEX IF EXISTS idx_toko_status_updated ON toko;