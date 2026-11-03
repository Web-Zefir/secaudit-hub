export type Severity = "CRITICAL" | "HIGH" | "MEDIUM" | "LOW";
export type FindingStatus = "OPEN" | "RESOLVED" | "FALSE_POSITIVE";
export type ScanStatus = "PENDING" | "RUNNING" | "COMPLETED" | "FAILED";
export type TargetSource = "DEMO_PRESET" | "CUSTOM_CODE" | "GIT_REPO";

export interface User {
  id: string;
  email: string;
  name: string;
  role: string;
  createdAt: string;
}

export interface Project {
  id: string;
  name: string;
  slug: string;
  description: string;
  repoUrl: string;
  language: string;
  lastScanAt: string | null;
  lastSecurityScore: number | null;
  createdAt: string;
}

export interface ScanJob {
  id: string;
  projectId: string;
  scannerType: string;
  targetSource: TargetSource;
  targetRef: string;
  status: ScanStatus;
  securityScore: number;
  filesScanned: number;
  durationMs: number;
  criticalCount: number;
  highCount: number;
  mediumCount: number;
  lowCount: number;
  errorMessage: string;
  startedAt: string | null;
  finishedAt: string | null;
  createdAt: string;
}

export interface Finding {
  id: string;
  scanJobId: string;
  projectId: string;
  ruleId: string;
  category: string;
  severity: Severity;
  title: string;
  description: string;
  filePath: string;
  lineNumber: number;
  codeSnippet: string;
  vulnerableLine: string;
  remediationAdvice: string;
  fixedCodeSnippet: string;
  status: FindingStatus;
  createdAt: string;
}

export interface CategoryStat {
  category: string;
  count: number;
}

export interface Dashboard {
  averageScore: number | null;
  projectCount: number;
  openCritical: number;
  openHigh: number;
  categories: CategoryStat[];
  projects: Project[];
}

export interface ScanResponse {
  scan: ScanJob;
  findings: Finding[];
}

export interface ModuleInfo {
  id: string;
  name: string;
  status: "available" | "roadmap";
}

export interface CodeFile {
  name: string;
  content: string;
}
