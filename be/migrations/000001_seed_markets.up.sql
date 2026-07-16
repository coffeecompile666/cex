-- Seed markets data into 'markets' table
INSERT INTO markets (name, code, created_at, updated_at) VALUES
('US Dollar', 'USD', NOW(), NOW()),
('Bitcoin', 'BTC', NOW(), NOW()),
('Ethereum', 'ETH', NOW(), NOW()),
('Solana', 'SOL', NOW(), NOW()),
('Dogecoin', 'DOGE', NOW(), NOW()),
('Cardano', 'ADA', NOW(), NOW())
ON CONFLICT (code) DO NOTHING;
