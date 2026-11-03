import { cx } from "../../lib/format";
import styles from "./Report.module.scss";

export function ScoreRing({ score }: { score: number }) {
  const tone = score >= 80 ? styles.good : score >= 50 ? styles.warn : styles.bad;
  return (
    <div className={styles.ringWrap}>
      <svg viewBox="0 0 36 36" className={styles.ring} aria-hidden="true">
        <circle className={styles.track} cx="18" cy="18" r="15.5" pathLength={100} />
        <circle className={cx(styles.arc, tone)} cx="18" cy="18" r="15.5" pathLength={100} strokeDasharray={`${score} 100`} />
      </svg>
      <div className={styles.ringLabel}>
        <strong>{score}</strong>
        <span>score</span>
      </div>
    </div>
  );
}
