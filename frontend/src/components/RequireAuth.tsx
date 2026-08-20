import { useEffect } from "react";
import { Navigate, Outlet } from "react-router";
import { useAuthStore } from "@/stores/auth";

function RequireAuth() {
  const { user, initialized, fetchUser } = useAuthStore();

  useEffect(() => {
    if (!initialized) {
      fetchUser();
    }
  }, [initialized, fetchUser]);

  if (!initialized) {
    return (
      <main className="flex min-h-screen items-center justify-center">
        <p className="text-muted-foreground">Loading…</p>
      </main>
    );
  }

  if (!user) {
    return <Navigate to="/login" replace />;
  }

  return <Outlet />;
}

export default RequireAuth;
