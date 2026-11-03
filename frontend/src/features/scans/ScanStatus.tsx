import { statusLabel } from "../../lib/format";
import type { ScanStatus } from "../../types/models";
import styles from "./Project.module.scss";

export function ScanStatusPill({ status }: { status: ScanStatus }) {
  return (
    <span className={styles.status}>
      <span className={`${styles.dot} ${styles[status]}`} />
      {statusLabel(status)}
    </span>
  );
}
