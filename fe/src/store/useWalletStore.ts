import { create } from 'zustand';
import { Wallet } from '../features/wallet/components/BalanceCard';

interface WalletState {
  wallet: Wallet | null;
  isLoading: boolean;
  error: string | null;
  fetchBalance: (userId: number) => Promise<void>;
  clearWallet: () => void;
}

export const useWalletStore = create<WalletState>((set) => ({
  wallet: null,
  isLoading: false,
  error: null,
  
  fetchBalance: async (userId: number) => {
    set({ isLoading: true, error: null });
    try {
      const response = await fetch(`/api/v1/wallets/${userId}`);
      const data = await response.json();
      
      if (!response.ok) {
        throw new Error(data.error || 'Failed to fetch balance');
      }

      set({ wallet: data as Wallet, isLoading: false });
    } catch (err: any) {
      set({ error: err.message || 'An error occurred', isLoading: false });
      throw err;
    }
  },

  clearWallet: () => set({ wallet: null, error: null }),
}));
