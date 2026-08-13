import { create } from "zustand";

interface DashboardState {
  status: string | null;
  loading: boolean;
  error: string | null;
  fetchStatus: () => Promise<void>;
}

export const useDashboardStore = create<DashboardState>((set) => ({
  status: null,
  loading: false,
  error: null,
  fetchStatus: async () => {
    set({ loading: true, error: null });
    try {
      const res = await fetch("/api/health");
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      const data = (await res.json()) as { status: string };
      set({ status: data.status, loading: false });
    } catch (e) {
      set({ error: e instanceof Error ? e.message : "Request failed", loading: false });
    }
  },
}));