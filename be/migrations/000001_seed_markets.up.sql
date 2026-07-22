-- Seed markets data into 'markets' table
INSERT INTO markets (name, symbol, decimals, precision, smallest_unit, is_base_currency, created_at, updated_at) VALUES
('US Dollar',  'USD',  2,  2,  'cent',       true,  NOW(), NOW()),
('Bitcoin',    'BTC',  8,  8,  'satoshi',    false, NOW(), NOW()),
('Ethereum',   'ETH',  18, 8,  'wei',        false, NOW(), NOW()),
('Solana',     'SOL',  9,  6,  'lamport',    false, NOW(), NOW()),
('Dogecoin',   'DOGE', 8,  4,  'koinu',      false, NOW(), NOW()),
('Cardano',    'ADA',  6,  6,  'lovelace',   false, NOW(), NOW())
ON CONFLICT (symbol) DO NOTHING;
