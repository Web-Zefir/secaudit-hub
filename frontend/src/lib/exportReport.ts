import type { Finding, ScanJob } from "../types/models";
import { statusLabel } from "./format";

export function downloadText(filename: string, content: string, mime: string): void {
  const blob = new Blob([content], { type: mime });
  const url = URL.createObjectURL(blob);
  const link = document.createElement("a");
  link.href = url;
  link.download = filename;
  link.click();
  URL.revokeObjectURL(url);
}

export function exportJSON(scan: ScanJob, findings: Finding[]): void {
  const payload = { scan, findings };
  downloadText(`secaudit-${scan.id}.json`, JSON.stringify(payload, null, 2), "application/json");
}

export function exportMarkdown(scan: ScanJob, findings: Finding[]): void {
  const lines = [
    `# SAST report ${scan.id}`,
    "",
    `- Status: ${scan.status}`,
    `- Security score: ${scan.securityScore}`,
    `- Files scanned: ${scan.filesScanned}`,
    `- Duration: ${scan.durationMs} ms`,
    `- Target: ${scan.targetSource} ${scan.targetRef}`,
    `- Counts: CRITICAL ${scan.criticalCount}, HIGH ${scan.highCount}, MEDIUM ${scan.mediumCount}, LOW ${scan.lowCount}`,
    "",
  ];
  findings.forEach((finding, index) => {
    lines.push(`## ${index + 1}. ${finding.title}`);
    lines.push("");
    lines.push(`- Rule: \`${finding.ruleId}\``);
    lines.push(`- Severity: ${finding.severity}`);
    lines.push(`- Category: ${finding.category}`);
    lines.push(`- Status: ${statusLabel(finding.status)}`);
    lines.push(`- Location: \`${finding.filePath}:${finding.lineNumber}\``);
    lines.push("");
    lines.push(finding.description);
    lines.push("");
    lines.push("```");
    lines.push(finding.codeSnippet);
    lines.push("```");
    lines.push("");
    lines.push(finding.remediationAdvice);
    lines.push("");
    lines.push("```");
    lines.push(finding.fixedCodeSnippet);
    lines.push("```");
    lines.push("");
  });
  downloadText(`secaudit-${scan.id}.md`, lines.join("\n"), "text/markdown");
}
