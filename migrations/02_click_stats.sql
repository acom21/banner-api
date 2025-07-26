-- Create click_stats table for tracking banner clicks
CREATE TABLE IF NOT EXISTS click_stats (
    banner_id INTEGER NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    count INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (banner_id, timestamp)
);

-- Create foreign key constraint to banners table
ALTER TABLE click_stats
ADD CONSTRAINT fk_click_stats_banner_id
FOREIGN KEY (banner_id) REFERENCES banners(id) ON DELETE CASCADE;

-- Create index for efficient time-based queries
CREATE INDEX IF NOT EXISTS idx_click_stats_timestamp
ON click_stats(timestamp);

-- Create index for banner + time range queries
CREATE INDEX IF NOT EXISTS idx_click_stats_banner_time
ON click_stats(banner_id, timestamp);
