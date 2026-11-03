import { Link, useNavigate, useParams } from "@tanstack/react-router";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import { api, errorMessage } from "../api/client";
import { Button } from "../components/ui/Button";
import { LaunchScanDialog } from "../features/scans/LaunchScanDialog";
import styles from "../features/scans/Project.module.scss";
import { ScanStatusPill } from "../features/scans/ScanStatus";
import { formatDuration, formatWhen, targetLabel } from "../lib/format";
import type { CodeFile, Project, ScanJob, TargetSource } from "../types/models";

export function ProjectPage() {
  const { projectId } = useParams({ from: "/app/projects/$projectId" });
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [open, setOpen] = useState(false);
  const [error, setError] = useState("");

  const project = useQuery({
    queryKey: ["project", projectId],
    queryFn: async () => (await api.get<Project>(`/projects/${projectId}`)).data,
  });
  const scans = useQuery({
    queryKey: ["scans", projectId],
    queryFn: async () => (await api.get<ScanJob[]>(`/projects/${projectId}/scans`)).data,
    refetchInterval: (query) => {
      const items = query.state.data;
      if (items?.some((item) => item.status === "PENDING" || item.status === "RUNNING")) return 1200;
      return false;
    },
  });
  const start = useMutation({
    mutationFn: async (input: { targetSource: TargetSource; repoUrl?: string; files?: CodeFile[] }) =>
      (await api.post<ScanJob>(`/projects/${projectId}/scans`, input)).data,
    onSuccess: async (job) => {
      setOpen(false);
      setError("");
      await queryClient.invalidateQueries({ queryKey: ["scans", projectId] });
      void navigate({ to: "/scans/$scanId", params: { scanId: job.id } });
    },
    onError: (err) => setError(errorMessage(err)),
  });

  return (
    <section>
      <nav className={styles.crumbs}>
        <Link to="/">Обзор</Link>
        <span>/</span>
        <span>{project.data?.name ?? "Проект"}</span>
      </nav>
      <header className={styles.header}>
        <div>
          <h1 className={styles.title}>{project.data?.name ?? "Проект"}</h1>
          <p className={styles.meta}>
            {project.data?.language ?? ""} · {project.data?.description || "Без описания"}
          </p>
        </div>
        <Button variant="primary" onClick={() => setOpen(true)}>
          Запустить SAST-аудит
        </Button>
      </header>
      <div className={styles.panel}>
        {scans.isLoading ? <div className={styles.empty}>Загрузка истории…</div> : null}
        {scans.data && scans.data.length === 0 ? <div className={styles.empty}>Сканов пока нет.</div> : null}
        {scans.data && scans.data.length > 0 ? (
          <table className={styles.table}>
            <thead>
              <tr>
                <th>Старт</th>
                <th>Источник</th>
                <th>Статус</th>
                <th>Score</th>
                <th>C / H / M / L</th>
                <th>Файлы</th>
                <th>Время</th>
              </tr>
            </thead>
            <tbody>
              {scans.data.map((scan) => (
                <tr key={scan.id}>
                  <td>
                    <Link to="/scans/$scanId" params={{ scanId: scan.id }}>
                      {formatWhen(scan.createdAt)}
                    </Link>
                  </td>
                  <td>
                    <div>{targetLabel(scan.targetSource)}</div>
                    <div className={styles.mono}>{scan.targetRef}</div>
                  </td>
                  <td>
                    <ScanStatusPill status={scan.status} />
                  </td>
                  <td>{scan.status === "COMPLETED" ? scan.securityScore : "—"}</td>
                  <td className={styles.mono}>
                    {scan.criticalCount}/{scan.highCount}/{scan.mediumCount}/{scan.lowCount}
                  </td>
                  <td>{scan.filesScanned}</td>
                  <td>{formatDuration(scan.durationMs)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        ) : null}
      </div>
      <LaunchScanDialog open={open} pending={start.isPending} error={error} onOpenChange={setOpen} onStart={(input) => start.mutate(input)} />
    </section>
  );
}
