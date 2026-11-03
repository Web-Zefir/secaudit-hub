import { Link, Outlet, useNavigate } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import * as Tooltip from "@radix-ui/react-tooltip";
import { useEffect } from "react";
import { api } from "../../api/client";
import { useAuth } from "../../context/AuthContext";
import type { ModuleInfo } from "../../types/models";
import { RoadmapBadge } from "../ui/Badge";
import { Button } from "../ui/Button";
import { Hint } from "../ui/Hint";
import styles from "./AppShell.module.scss";

export function AppShell() {
  const { user, loading, logout } = useAuth();
  const navigate = useNavigate();
  const modules = useQuery({
    queryKey: ["modules"],
    queryFn: async () => (await api.get<ModuleInfo[]>("/modules")).data,
  });

  useEffect(() => {
    if (!loading && !user) {
      void navigate({ to: "/login" });
    }
  }, [loading, user, navigate]);

  return (
    <Tooltip.Provider delayDuration={200}>
      <div className={styles.shell}>
        <aside className={styles.sidebar}>
          <div className={styles.brand}>
            <span className={styles.mark}>SA</span>
            <span className={styles.word}>SecAudit Hub</span>
          </div>
          <nav className={styles.nav}>
            <Link to="/" activeOptions={{ exact: true }} className={styles.item} activeProps={{ className: styles.active }}>
              Обзор
            </Link>
            <div className={styles.section}>Модули</div>
            {(modules.data ?? fallbackModules).map((module) =>
              module.status === "available" ? (
                <Link key={module.id} to="/" className={styles.item}>
                  <span>{module.name}</span>
                </Link>
              ) : (
                <Hint key={module.id} label="Модуль запланирован и не выполняет проверки">
                  <span className={styles.disabled}>
                    <span>{module.name}</span>
                    <RoadmapBadge />
                  </span>
                </Hint>
              ),
            )}
          </nav>
          <div className={styles.user}>
            {loading || !user ? <div className={styles.skeleton} /> : <div className={styles.email}>{user.email}</div>}
            <Button onClick={() => void logout()}>Выйти</Button>
          </div>
        </aside>
        <main className={styles.main}>{loading || !user ? <div className={styles.skeleton} /> : <Outlet />}</main>
      </div>
    </Tooltip.Provider>
  );
}

const fallbackModules: ModuleInfo[] = [
  { id: "sast", name: "SAST", status: "available" },
  { id: "dast", name: "DAST", status: "roadmap" },
  { id: "iac", name: "IaC", status: "roadmap" },
  { id: "dbsec", name: "DB Security", status: "roadmap" },
];
