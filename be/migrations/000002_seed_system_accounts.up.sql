-- Seed system accounts for each market (asset)
-- These accounts have user_id = NULL and is_system_account = true.
-- The type is set to 'MAIN_POOL' representing the main system liquidity pool/vault for that asset.

INSERT INTO accounts (user_id, is_system_account, market_id, type, created_at, updated_at)
SELECT 
    NULL AS user_id, 
    true AS is_system_account, 
    id AS market_id, 
    'MAIN_POOL' AS type, 
    NOW() AS created_at, 
    NOW() AS updated_at
FROM markets
ON CONFLICT (user_id, market_id, type) DO NOTHING;
