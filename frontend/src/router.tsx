import { createBrowserRouter } from "react-router";
import Dashboard from "./components/Dashboard";
import Login from "./components/Login";
import Register from "./components/Register";
import RequireAuth from "./components/RequireAuth";

export const router = createBrowserRouter([
  {
    path: "/login",
    element: <Login />,
  },
  {
    path: "/register",
    element: <Register />,
  },
  {
    element: <RequireAuth />,
    children: [
      {
        path: "/",
        element: <Dashboard />,
      },
    ],
  },
]);
