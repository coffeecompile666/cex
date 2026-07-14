import React, { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import BalanceCard from '../features/wallet/components/BalanceCard';
import { useUserStore } from '../store/useUserStore';

export default function Dashboard() {
  const { currentUser } = useUserStore();
  const navigate = useNavigate();

  // Protect route: redirect to login if no user
  useEffect(() => {
    if (!currentUser) {
      navigate('/login');
    }
  }, [currentUser, navigate]);

  if (!currentUser) return null; // Avoid flashing content before redirect

  return (
    <div className="container">
      <div className="grid-2">
        <div>
          <BalanceCard />
        </div>
        
        <div>
          {/* We will add OrderBook or Trade UI here later */}
          <div className="card" style={{ height: '100%', minHeight: '250px', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
            <p className="text-muted">Trade Module Coming Soon...</p>
          </div>
        </div>
      </div>
    </div>
  );
}
