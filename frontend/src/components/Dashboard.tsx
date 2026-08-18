import { useEffect } from "react";
import { Link } from "react-router";
import { useDashboardStore } from "../stores/dashboard";

function Dashboard() {
  const { status, loading, error, fetchStatus } = useDashboardStore();

  useEffect(() => {
    fetchStatus();
  }, [fetchStatus]);

  return (
    <main>
      <h1 className="font-bold text-2xl">Physiq</h1>
      <p>Weight and body measurement tracking.</p>
      {loading && <p>Loading…</p>}
      {error && <p className="text-red-500">{error}</p>}
      {status && <p>API status: {status}</p>}
      <Link to="/register" className="text-primary underline underline-offset-4">
        Create an account
      </Link>
    </main>
  );
}

export default Dashboard;
