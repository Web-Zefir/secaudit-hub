import { useNavigate } from "@tanstack/react-router";
import { createContext, useContext, useEffect, useMemo, useState, type ReactNode } from "react";
import { api, setUnauthorizedHandler } from "../api/client";
import type { User } from "../types/models";

interface AuthContextValue {
  user: User | null;
  loading: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (name: string, email: string, password: string) => Promise<void>;
  loginDemo: () => Promise<void>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const navigate = useNavigate();
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let live = true;
    api
      .get<User>("/auth/me")
      .then((response) => {
        if (live) setUser(response.data);
      })
      .catch(() => {
        if (live) setUser(null);
      })
      .finally(() => {
        if (live) setLoading(false);
      });
    return () => {
      live = false;
    };
  }, []);

  useEffect(() => {
    setUnauthorizedHandler(() => {
      setUser(null);
      void navigate({ to: "/login" });
    });
  }, [navigate]);

  const value = useMemo<AuthContextValue>(
    () => ({
      user,
      loading,
      async login(email, password) {
        const response = await api.post<User>("/auth/login", { email, password });
        setUser(response.data);
      },
      async register(name, email, password) {
        const response = await api.post<User>("/auth/register", { name, email, password });
        setUser(response.data);
      },
      async loginDemo() {
        const response = await api.post<User>("/auth/demo");
        setUser(response.data);
      },
      async logout() {
        await api.post("/auth/logout");
        setUser(null);
        void navigate({ to: "/login" });
      },
    }),
    [user, loading, navigate],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth(): AuthContextValue {
  const value = useContext(AuthContext);
  if (!value) {
    throw new Error("useAuth must be used inside AuthProvider");
  }
  return value;
}
