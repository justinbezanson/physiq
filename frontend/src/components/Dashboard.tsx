import { useEffect } from "react";
import { useNavigate } from "react-router";
import { useDashboardStore } from "../stores/dashboard";
import { useAuthStore } from "@/stores/auth";
import { Button } from "@/components/ui/button";

function Dashboard() {
  const navigate = useNavigate();
  const { status, loading, error, fetchStatus } = useDashboardStore();
  const { user, logout } = useAuthStore();

  useEffect(() => {
    fetchStatus();
  }, [fetchStatus]);

  async function handleLogout() {
    await logout();
    navigate("/login");
  }

  return (
    <main className="min-h-screen p-8">
      <div className="flex items-center justify-between">
        <h1 className="font-bold text-2xl">Physiq</h1>
        <div className="flex items-center gap-3">
          <span className="text-sm text-muted-foreground">{user?.name}</span>
          <Button variant="outline" size="sm" onClick={handleLogout}>
            Log out
          </Button>
        </div>
      </div>
      <p>Weight and body measurement tracking.</p>
      {loading && <p>Loading…</p>}
      {error && <p className="text-red-500">{error}</p>}
      {status && <p>API status: {status}</p>}
    </main>
  );
}

export default Dashboard;
