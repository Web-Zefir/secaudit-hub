import { Link } from "@tanstack/react-router";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import { api, errorMessage } from "../api/client";
import { Button } from "../components/ui/Button";
import { CreateProjectDialog } from "../features/dashboard/CreateProjectDialog";
import styles from "../features/dashboard/Dashboard.module.scss";
import { formatScore, formatWhen } from "../lib/format";
import type { Dashboard, Project } from "../types/models";

export function DashboardPage() {
  const queryClient = useQueryClient();
  const [open, setOpen] = useState(false);
  const [error, setError] = useState("");
  const dashboard = useQuery({
    queryKey: ["dashboard"],
    queryFn: async () => (await api.get<Dashboard>("/dashboard")).data,
  });
  const create = useMutation({
    mutationFn: async (input: { name: string; description: string; repoUrl: string; language: string }) =>
      (await api.post<Project>("/projects", input)).data,
    onSuccess: async () => {
      setOpen(false);
      setError("");
      await queryClient.invalidateQueries({ queryKey: ["dashboard"] });
    },
    onError: (err) => setError(errorMessage(err)),
  });

  const data = dashboard.data;
  const maxCategory = Math.max(1, ...(data?.categories.map((item) => item.count) ?? [1]));

  return (
    <section>
      <header className={styles.header}>
        <div>
          <h1 className={styles.title}>Обзор</h1>
          <p className={styles.subtitle}>Сводка по проектам и открытым находкам SAST.</p>
        </div>
        <Button variant="primary" onClick={() => setOpen(true)}>
          Новый проект
        </Button>
      </header>
      <div className={styles.metrics}>
        <article className={styles.metric}>
          <div className={styles.label}>Средний Security Score</div>
          <div className={styles.value}>{data ? formatScore(data.averageScore) : "…"}</div>
        </article>
        <article className={styles.metric}>
          <div className={styles.label}>Проекты</div>
          <div className={styles.value}>{data?.projectCount ?? "…"}</div>
        </article>
        <article className={styles.metric}>
          <div className={styles.label}>Открытые Critical</div>
          <div className={`${styles.value} ${styles.danger}`}>{data?.openCritical ?? "…"}</div>
        </article>
        <article className={styles.metric}>
          <div className={styles.label}>Открытые High</div>
          <div className={styles.value}>{data?.openHigh ?? "…"}</div>
        </article>
      </div>
      <div className={styles.layout}>
        <div className={styles.panel}>
          <div className={styles.panelHead}>Проекты</div>
          {dashboard.isLoading ? <div className={styles.empty}>Загрузка…</div> : null}
          {data && data.projects.length === 0 ? (
            <div className={styles.empty}>Проектов пока нет. Создайте первый и запустите SAST.</div>
          ) : null}
          {data && data.projects.length > 0 ? (
            <table className={styles.table}>
              <thead>
                <tr>
                  <th>Проект</th>
                  <th>Язык</th>
                  <th>Score</th>
                  <th>Последний скан</th>
                </tr>
              </thead>
              <tbody>
                {data.projects.map((project) => (
                  <tr key={project.id}>
                    <td>
                      <Link to="/projects/$projectId" params={{ projectId: project.id }} className={styles.name}>
                        {project.name}
                      </Link>
                      <div className={styles.meta}>{project.slug}</div>
                    </td>
                    <td>{project.language}</td>
                    <td>{formatScore(project.lastSecurityScore)}</td>
                    <td>{formatWhen(project.lastScanAt)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          ) : null}
        </div>
        <aside className={styles.panel}>
          <div className={styles.panelHead}>Категории OWASP</div>
          {data && data.categories.length === 0 ? <div className={styles.empty}>Нет открытых находок.</div> : null}
          <div className={styles.bars}>
            {data?.categories.map((item) => (
              <div key={item.category}>
                <div className={styles.barRow}>
                  <span className={styles.barLabel}>{item.category}</span>
                  <span>{item.count}</span>
                </div>
                <svg className={styles.trackSvg} viewBox="0 0 100 4" preserveAspectRatio="none" aria-hidden="true">
                  <rect className={styles.trackBg} width={100} height={4} />
                  <rect className={styles.fill} width={Math.max(4, (item.count / maxCategory) * 100)} height={4} />
                </svg>
              </div>
            ))}
          </div>
        </aside>
      </div>
      <CreateProjectDialog
        open={open}
        pending={create.isPending}
        error={error}
        onOpenChange={setOpen}
        onSubmit={(input) => create.mutate(input)}
      />
    </section>
  );
}
