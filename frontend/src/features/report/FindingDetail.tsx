import { useState } from "react";
import { NeutralBadge, SeverityBadge } from "../../components/ui/Badge";
import { Button } from "../../components/ui/Button";
import { statusLabel } from "../../lib/format";
import type { Finding, FindingStatus } from "../../types/models";
import styles from "./Report.module.scss";

interface FindingDetailProps {
  finding: Finding;
  pending: boolean;
  onStatus: (status: FindingStatus) => void;
}

export function FindingDetail({ finding, pending, onStatus }: FindingDetailProps) {
  const [copied, setCopied] = useState(false);

  async function copyFix() {
    await navigator.clipboard.writeText(finding.fixedCodeSnippet);
    setCopied(true);
    window.setTimeout(() => setCopied(false), 1200);
  }

  return (
    <article className={styles.detail}>
      <h2 className={styles.title}>{finding.title}</h2>
      <div className={styles.tags}>
        <SeverityBadge severity={finding.severity} />
        <NeutralBadge>{finding.ruleId}</NeutralBadge>
        <NeutralBadge>{finding.category}</NeutralBadge>
        <NeutralBadge>{statusLabel(finding.status)}</NeutralBadge>
      </div>
      <div className={styles.mono}>
        {finding.filePath}:{finding.lineNumber}
      </div>
      <p className={styles.advice}>{finding.description}</p>
      <pre className={styles.code}>
        <code>{finding.codeSnippet}</code>
      </pre>
      <div className={styles.diff}>
        <div>
          <h3>Уязвимый код</h3>
          <pre>
            <code>{finding.vulnerableLine}</code>
          </pre>
        </div>
        <div>
          <h3>Исправленный код</h3>
          <pre>
            <code>{finding.fixedCodeSnippet}</code>
          </pre>
        </div>
      </div>
      <p className={styles.advice}>{finding.remediationAdvice}</p>
      <div className={styles.row}>
        <Button onClick={() => void copyFix()}>{copied ? "Скопировано" : "Скопировать исправленный код"}</Button>
        <Button variant="primary" disabled={pending} onClick={() => onStatus("RESOLVED")}>
          Mark as Resolved
        </Button>
        <Button disabled={pending} onClick={() => onStatus("FALSE_POSITIVE")}>
          False Positive
        </Button>
        {finding.status !== "OPEN" ? (
          <Button disabled={pending} onClick={() => onStatus("OPEN")}>
            Вернуть в Open
          </Button>
        ) : null}
      </div>
    </article>
  );
}
