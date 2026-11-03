import { Link, useParams } from "@tanstack/react-router";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import * as ScrollArea from "@radix-ui/react-scroll-area";
import { useMemo, useState } from "react";
import { api } from "../api/client";
import { SeverityBadge } from "../components/ui/Badge";
import { Button } from "../components/ui/Button";
import { SelectInput, TextInput } from "../components/ui/Field";
import { FindingDetail } from "../features/report/FindingDetail";
import styles from "../features/report/Report.module.scss";
import { ScoreRing } from "../features/report/ScoreRing";
import { ScanStatusPill } from "../features/scans/ScanStatus";
import { exportJSON, exportMarkdown } from "../lib/exportReport";
import { formatDuration, statusLabel } from "../lib/format";
import type { Finding, FindingStatus, ScanResponse, Severity } from "../types/models";

export function ScanPage() {
  const { scanId } = useParams({ from: "/app/scans/$scanId" });
  const queryClient = useQueryClient();
  const [query, setQuery] = useState("");
  const [severity, setSeverity] = useState<Severity | "ALL">("ALL");
  const [category, setCategory] = useState("ALL");
  const [status, setStatus] = useState<FindingStatus | "ALL">("ALL");
  const [selectedId, setSelectedId] = useState<string | null>(null);

  const report = useQuery({
    queryKey: ["scan", scanId],
    queryFn: async () => (await api.get<ScanResponse>(`/scans/${scanId}`)).data,
    refetchInterval: (current) => {
      const value = current.state.data?.scan.status;
      return value === "PENDING" || value === "RUNNING" ? 1200 : false;
    },
  });
  const update = useMutation({
    mutationFn: async (input: { id: string; status: FindingStatus }) =>
      (await api.patch<Finding>(`/findings/${input.id}/status`, { status: input.status })).data,
    onSuccess: async (finding) => {
      queryClient.setQueryData<ScanResponse>(["scan", scanId], (old) => {
        if (!old) return old;
        return { ...old, findings: old.findings.map((item) => (item.id === finding.id ? finding : item)) };
      });
      await queryClient.invalidateQueries({ queryKey: ["dashboard"] });
    },
  });

  const findings = report.data?.findings ?? [];
  const categories = useMemo(() => Array.from(new Set(findings.map((item) => item.category))), [findings]);
  const filtered = findings.filter((item) => {
    const haystack = `${item.title} ${item.filePath} ${item.ruleId} ${item.description} ${item.category}`.toLowerCase();
    if (query && !haystack.includes(query.toLowerCase())) return false;
    if (severity !== "ALL" && item.severity !== severity) return false;
    if (category !== "ALL" && item.category !== category) return false;
    if (status !== "ALL" && item.status !== status) return false;
    return true;
  });
  const selected = filtered.find((item) => item.id === selectedId) ?? filtered[0] ?? null;
  const scan = report.data?.scan;

  return (
    <section>
      <nav className={styles.crumbs}>
        <Link to="/">Обзор</Link>
        <span>/</span>
        {scan ? <Link to="/projects/$projectId" params={{ projectId: scan.projectId }}>Проект</Link> : <span>Проект</span>}
        <span>/</span>
        <span>Отчёт</span>
      </nav>
      <header className={styles.header}>
        <div>
          <h1 className={styles.title}>SAST-отчёт</h1>
          {scan ? <ScanStatusPill status={scan.status} /> : null}
        </div>
        <div className={styles.actions}>
          <Button disabled={!scan} onClick={() => scan && exportJSON(scan, findings)}>
            Экспорт JSON
          </Button>
          <Button disabled={!scan} onClick={() => scan && exportMarkdown(scan, findings)}>
            Экспорт Markdown
          </Button>
        </div>
      </header>
      {scan?.status === "FAILED" ? <div className={`${styles.banner} ${styles.failed}`}>{scan.errorMessage}</div> : null}
      {scan && (scan.status === "PENDING" || scan.status === "RUNNING") ? (
        <div className={styles.banner}>Скан выполняется. Статус обновляется автоматически.</div>
      ) : null}
      {scan?.status === "COMPLETED" ? (
        <div className={styles.summary}>
          <article className={styles.card}>
            <ScoreRing score={scan.securityScore} />
          </article>
          <article className={styles.card}>
            <div className={styles.label}>Время</div>
            <div className={styles.figure}>{formatDuration(scan.durationMs)}</div>
          </article>
          <article className={styles.card}>
            <div className={styles.label}>Файлы</div>
            <div className={styles.figure}>{scan.filesScanned}</div>
          </article>
          <Count label="Critical" value={scan.criticalCount} />
          <Count label="High" value={scan.highCount} />
          <Count label="Medium / Low" value={`${scan.mediumCount} / ${scan.lowCount}`} />
        </div>
      ) : null}
      <div className={styles.filters}>
        <TextInput value={query} onChange={(event) => setQuery(event.target.value)} placeholder="Поиск по файлу, правилу, тексту" />
        <SelectInput value={severity} onChange={(event) => setSeverity(event.target.value as Severity | "ALL")}>
          <option value="ALL">Все severity</option>
          <option value="CRITICAL">Critical</option>
          <option value="HIGH">High</option>
          <option value="MEDIUM">Medium</option>
          <option value="LOW">Low</option>
        </SelectInput>
        <SelectInput value={category} onChange={(event) => setCategory(event.target.value)}>
          <option value="ALL">Все категории</option>
          {categories.map((item) => (
            <option key={item} value={item}>
              {item}
            </option>
          ))}
        </SelectInput>
        <SelectInput value={status} onChange={(event) => setStatus(event.target.value as FindingStatus | "ALL")}>
          <option value="ALL">Все статусы</option>
          <option value="OPEN">{statusLabel("OPEN")}</option>
          <option value="RESOLVED">{statusLabel("RESOLVED")}</option>
          <option value="FALSE_POSITIVE">{statusLabel("FALSE_POSITIVE")}</option>
        </SelectInput>
      </div>
      {scan?.status === "COMPLETED" && filtered.length === 0 ? <div className={styles.empty}>Нет находок по текущему фильтру.</div> : null}
      {filtered.length > 0 && selected ? (
        <div className={styles.split}>
          <ScrollArea.Root className={styles.list}>
            <ScrollArea.Viewport className={styles.viewport}>
              {filtered.map((finding) => (
                <button
                  key={finding.id}
                  type="button"
                  className={finding.id === selected.id ? `${styles.item} ${styles.itemActive}` : styles.item}
                  onClick={() => setSelectedId(finding.id)}
                >
                  <span className={styles.itemTitle}>{finding.title}</span>
                  <span className={styles.path}>
                    {finding.filePath}:{finding.lineNumber}
                  </span>
                  <SeverityBadge severity={finding.severity} />
                </button>
              ))}
            </ScrollArea.Viewport>
            <ScrollArea.Scrollbar orientation="vertical">
              <ScrollArea.Thumb />
            </ScrollArea.Scrollbar>
          </ScrollArea.Root>
          <FindingDetail finding={selected} pending={update.isPending} onStatus={(next) => update.mutate({ id: selected.id, status: next })} />
        </div>
      ) : null}
    </section>
  );
}

function Count({ label, value }: { label: string; value: number | string }) {
  return (
    <article className={styles.card}>
      <div className={styles.label}>{label}</div>
      <div className={styles.figure}>{value}</div>
    </article>
  );
}
