import React, { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useUserStore } from '../store/useUserStore';
import LoginForm from '../features/user/components/LoginForm';

export default function Login() {
  const { currentUser } = useUserStore();
  const navigate = useNavigate();

  // If already logged in, redirect to dashboard
  useEffect(() => {
    if (currentUser) {
      navigate('/');
    }
  }, [currentUser, navigate]);

  return (
    <div className="container" style={{ maxWidth: '600px' }}>
      <LoginForm />
    </div>
  );
}
