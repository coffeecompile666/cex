import React, { useEffect } from 'react';
import { useUserStore } from '../../../store/useUserStore';
import { useWalletStore } from '../../../store/useWalletStore';
import { Button } from '../../../components/ui/Button';
import { Card } from '../../../components/ui/Card';

export interface Wallet {
  id: number;
  user_id: number;
  balance: number;
}

export default function BalanceCard() {
  const { currentUser } = useUserStore();
  const { wallet, isLoading, error, fetchBalance } = useWalletStore();

  useEffect(() => {
    if (currentUser?.id) {
      fetchBalance(currentUser.id);
    }
  }, [currentUser, fetchBalance]);

  if (!currentUser) {
    return (
      <Card style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', height: '100%', minHeight: '250px' }}>
        <p className="text-muted">Please register or login to see your balance.</p>
      </Card>
    );
  }

  return (
    <Card title="Your Assets">
      {isLoading && !wallet ? (
        <p className="text-muted">Loading balance...</p>
      ) : wallet ? (
        <div>
          <div style={{ padding: '1.5rem', backgroundColor: '#f0fdf4', borderRadius: '8px', border: '1px solid #bbf7d0', marginBottom: '1rem' }}>
            <p className="text-muted" style={{ marginBottom: '0.5rem', color: '#166534' }}>Estimated Balance</p>
            <h3 style={{ fontSize: '2.5rem', color: '#15803d', margin: 0 }}>
              ${wallet.balance.toFixed(2)}
            </h3>
          </div>
          
          <Button style={{ backgroundColor: '#1f2937' }} onClick={() => fetchBalance(currentUser.id)} isLoading={isLoading}>
            Refresh Balance
          </Button>
          
          <p className="text-muted" style={{ marginTop: '1rem', fontSize: '0.8rem' }}>
            Wallet automatically created by the internal Wallet Module when you registered!
          </p>
        </div>
      ) : (
        <p style={{ color: 'red' }}>{error || 'Error loading wallet.'}</p>
      )}
    </Card>
  );
}
