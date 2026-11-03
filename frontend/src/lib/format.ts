import type { FindingStatus, ScanStatus, Severity, TargetSource } from "../types/models";

export function formatWhen(value: string | null): string {
  if (!value) return "—";
  return new Intl.DateTimeFormat("ru-RU", {
    dateStyle: "medium",
    timeStyle: "short",
  }).format(new Date(value));
}

export function formatDuration(ms: number): string {
  if (ms < 1000) return `${ms} мс`;
  return `${(ms / 1000).toFixed(1)} с`;
}

export function formatScore(score: number | null): string {
  if (score === null) return "—";
  return score.toFixed(0);
}

export function severityLabel(severity: Severity): string {
  switch (severity) {
    case "CRITICAL":
      return "Critical";
    case "HIGH":
      return "High";
    case "MEDIUM":
      return "Medium";
    case "LOW":
      return "Low";
  }
}

export function statusLabel(status: FindingStatus | ScanStatus): string {
  switch (status) {
    case "OPEN":
      return "Open";
    case "RESOLVED":
      return "Resolved";
    case "FALSE_POSITIVE":
      return "False Positive";
    case "PENDING":
      return "Pending";
    case "RUNNING":
      return "Running";
    case "COMPLETED":
      return "Completed";
    case "FAILED":
      return "Failed";
  }
}

export function targetLabel(source: TargetSource): string {
  switch (source) {
    case "DEMO_PRESET":
      return "Демо-стенд";
    case "CUSTOM_CODE":
      return "Свой код";
    case "GIT_REPO":
      return "Git";
  }
}

export function cx(...parts: Array<string | false | undefined>): string {
  return parts.filter((part): part is string => Boolean(part)).join(" ");
}
