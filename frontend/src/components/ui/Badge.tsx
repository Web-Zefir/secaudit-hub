import type { Severity } from "../../types/models";
import { severityLabel } from "../../lib/format";
import styles from "./Badge.module.scss";

export function SeverityBadge({ severity }: { severity: Severity }) {
  return <span className={`${styles.badge} ${styles[severity]}`}>{severityLabel(severity)}</span>;
}

export function NeutralBadge({ children }: { children: string }) {
  return <span className={`${styles.badge} ${styles.neutral}`}>{children}</span>;
}

export function RoadmapBadge() {
  return <span className={`${styles.badge} ${styles.roadmap}`}>In Roadmap</span>;
}
