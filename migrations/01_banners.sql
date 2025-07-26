-- Create banners table
CREATE TABLE IF NOT EXISTS banners (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index on name for faster lookups
CREATE INDEX IF NOT EXISTS idx_banners_name ON banners(name);

-- Insert some sample banners for testing
INSERT INTO banners (name) VALUES
    ('Banner 1'),
    ('Banner 2'),
    ('Banner 3'),
    ('Banner 4'),
    ('Banner 5'),
    ('Banner 6'),
    ('Banner 7'),
    ('Banner 8'),
    ('Banner 9'),
    ('Banner 10')
ON CONFLICT DO NOTHING;
