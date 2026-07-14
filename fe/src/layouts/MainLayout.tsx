import React from 'react';
import { Outlet, useNavigate } from 'react-router-dom';
import { useUserStore } from '../store/useUserStore';
import { Button } from '../components/ui/Button';

export default function MainLayout() {
  const { currentUser, logout } = useUserStore();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <div className="App">
      <header style={{ 
        backgroundColor: 'white', 
        padding: '1rem 2rem', 
        borderBottom: '1px solid var(--border-color)',
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center'
      }}>
        <h1 style={{ margin: 0, fontSize: '1.5rem', fontWeight: 'bold' }}>
          Icon <span className="text-crypto">Exchange</span>
        </h1>
        
        {currentUser && (
          <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
            <span className="text-muted">Welcome, {currentUser.username}</span>
            <div style={{ width: '80px' }}>
              <Button onClick={handleLogout} style={{ padding: '0.4rem 0.8rem', fontSize: '0.875rem' }}>Logout</Button>
            </div>
          </div>
        )}
      </header>
      <main style={{ padding: '2rem 0' }}>
        <Outlet />
      </main>
    </div>
  );
}
