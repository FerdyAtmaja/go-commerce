-- Update store status enum to include all required states
ALTER TABLE toko MODIFY COLUMN status ENUM('pending', 'active', 'inactive', 'suspended') DEFAULT 'pending';

-- Update existing stores to have proper status
-- Assuming current 'active' stores should remain 'active'
-- and 'suspended' stores should remain 'suspended'
UPDATE toko SET status = 'active' WHERE status = 'active';
UPDATE toko SET status = 'suspended' WHERE status = 'suspended';

-- Add index for better performance
CREATE INDEX IF NOT EXISTS idx_toko_status_updated ON toko(status);