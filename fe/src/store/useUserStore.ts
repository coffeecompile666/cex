import { create } from 'zustand';
import { User } from '../features/user/components/LoginForm';

interface UserState {
  currentUser: User | null;
  isLoading: boolean;
  error: string | null;
  login: (username: string, email: string) => Promise<void>;
  logout: () => void;
}

export const useUserStore = create<UserState>((set) => ({
  currentUser: null,
  isLoading: false,
  error: null,
  
  login: async (username: string, email: string) => {
    set({ isLoading: true, error: null });
    try {
      const response = await fetch('/api/v1/users/', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, email })
      });
      
      const data = await response.json();
      
      if (!response.ok) {
        throw new Error(data.error || 'Failed to create user');
      }

      set({ currentUser: data as User, isLoading: false });
    } catch (err: any) {
      set({ error: err.message || 'An error occurred', isLoading: false });
      throw err;
    }
  },

  logout: () => set({ currentUser: null, error: null }),
}));
