import React, { useState, FormEvent } from 'react';
import { useUserStore } from '../../../store/useUserStore';
import { Button } from '../../../components/ui/Button';
import { Input } from '../../../components/ui/Input';
import { Card } from '../../../components/ui/Card';

export interface User {
  id: number;
  username: string;
  email: string;
}

export default function LoginForm() {
  const [username, setUsername] = useState<string>('');
  const [email, setEmail] = useState<string>('');
  
  // Get state and actions from Zustand store
  const { login, isLoading, error } = useUserStore();

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    try {
      await login(username, email);
    } catch (err) {
      // Error is already handled and stored in Zustand
      console.error(err);
    }
  };

  return (
    <Card title="Welcome / Quick Start">
      <p className="text-muted" style={{ marginBottom: '1.5rem' }}>
        Register a new mock user to see the modular monolith in action.
      </p>
      
      <form onSubmit={handleSubmit}>
        <Input 
          label="Username"
          type="text" 
          placeholder="e.g. Satoshi" 
          value={username}
          onChange={e => setUsername(e.target.value)}
          required
        />
        
        <Input 
          label="Email"
          type="email" 
          placeholder="satoshi@crypto.com" 
          value={email}
          onChange={e => setEmail(e.target.value)}
          required
        />
        
        {error && <p style={{ color: 'red', fontSize: '0.875rem', marginBottom: '1rem' }}>{error}</p>}
        
        <Button type="submit" isLoading={isLoading}>
          Create Account
        </Button>
      </form>
    </Card>
  );
}
