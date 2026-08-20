import { create } from "zustand";

interface User {
  id: number;
  name: string;
  email: string;
}

interface AuthState {
  user: User | null;
  initialized: boolean;
  fetchUser: () => Promise<User | null>;
  logout: () => Promise<void>;
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  initialized: false,
  fetchUser: async () => {
    try {
      const res = await fetch("/api/me");
      if (!res.ok) {
        set({ user: null, initialized: true });
        return null;
      }
      const user = (await res.json()) as User;
      set({ user, initialized: true });
      return user;
    } catch {
      set({ user: null, initialized: true });
      return null;
    }
  },
  logout: async () => {
    try {
      await fetch("/api/auth/logout", { method: "POST" });
    } finally {
      set({ user: null });
    }
  },
}));
