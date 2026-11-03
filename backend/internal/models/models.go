package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	RoleAuditor = "auditor"

	ScannerSAST  = "SAST"
	ScannerDAST  = "DAST"
	ScannerIaC   = "IAC"
	ScannerDBSec = "DBSEC"

	TargetDemo   = "DEMO_PRESET"
	TargetCustom = "CUSTOM_CODE"
	TargetGit    = "GIT_REPO"

	ScanPending   = "PENDING"
	ScanRunning   = "RUNNING"
	ScanCompleted = "COMPLETED"
	ScanFailed    = "FAILED"

	SeverityCritical = "CRITICAL"
	SeverityHigh     = "HIGH"
	SeverityMedium   = "MEDIUM"
	SeverityLow      = "LOW"

	FindingOpen          = "OPEN"
	FindingResolved      = "RESOLVED"
	FindingFalsePositive = "FALSE_POSITIVE"
)

type User struct {
	ID           string    `gorm:"type:uuid;primaryKey" json:"id"`
	Email        string    `gorm:"uniqueIndex;size:255;not null" json:"email"`
	Name         string    `gorm:"size:120;not null" json:"name"`
	PasswordHash string    `gorm:"not null" json:"-"`
	Role         string    `gorm:"size:32;not null" json:"role"`
	CreatedAt    time.Time `json:"createdAt"`
}

func (u *User) BeforeCreate(_ *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.NewString()
	}
	return nil
}

type Project struct {
	ID                string     `gorm:"type:uuid;primaryKey" json:"id"`
	UserID            string     `gorm:"type:uuid;not null;index;uniqueIndex:idx_project_user_slug" json:"-"`
	Name              string     `gorm:"size:120;not null" json:"name"`
	Slug              string     `gorm:"size:140;not null;uniqueIndex:idx_project_user_slug" json:"slug"`
	Description       string     `gorm:"type:text" json:"description"`
	RepoURL           string     `gorm:"size:500" json:"repoUrl"`
	Language          string     `gorm:"size:40" json:"language"`
	LastScanAt        *time.Time `json:"lastScanAt"`
	LastSecurityScore *int       `json:"lastSecurityScore"`
	CreatedAt         time.Time  `json:"createdAt"`
}

func (p *Project) BeforeCreate(_ *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.NewString()
	}
	return nil
}

type ScanJob struct {
	ID            string     `gorm:"type:uuid;primaryKey" json:"id"`
	ProjectID     string     `gorm:"type:uuid;index;not null" json:"projectId"`
	ScannerType   string     `gorm:"size:16;not null;index" json:"scannerType"`
	TargetSource  string     `gorm:"size:32;not null" json:"targetSource"`
	TargetRef     string     `gorm:"size:500" json:"targetRef"`
	Status        string     `gorm:"size:16;not null;index" json:"status"`
	SecurityScore int        `json:"securityScore"`
	FilesScanned  int        `json:"filesScanned"`
	DurationMs    int64      `json:"durationMs"`
	CriticalCount int        `json:"criticalCount"`
	HighCount     int        `json:"highCount"`
	MediumCount   int        `json:"mediumCount"`
	LowCount      int        `json:"lowCount"`
	ErrorMessage  string     `gorm:"type:text" json:"errorMessage"`
	StartedAt     *time.Time `json:"startedAt"`
	FinishedAt    *time.Time `json:"finishedAt"`
	CreatedAt     time.Time  `json:"createdAt"`
}

func (s *ScanJob) BeforeCreate(_ *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.NewString()
	}
	return nil
}

type Finding struct {
	ID                string    `gorm:"type:uuid;primaryKey" json:"id"`
	ScanJobID         string    `gorm:"type:uuid;index;not null" json:"scanJobId"`
	ProjectID         string    `gorm:"type:uuid;index;not null" json:"projectId"`
	RuleID            string    `gorm:"size:32;not null;index" json:"ruleId"`
	Category          string    `gorm:"size:80;not null" json:"category"`
	Severity          string    `gorm:"size:16;not null;index" json:"severity"`
	Title             string    `gorm:"size:200;not null" json:"title"`
	Description       string    `gorm:"type:text" json:"description"`
	FilePath          string    `gorm:"size:500" json:"filePath"`
	LineNumber        int       `json:"lineNumber"`
	CodeSnippet       string    `gorm:"type:text" json:"codeSnippet"`
	VulnerableLine    string    `gorm:"type:text" json:"vulnerableLine"`
	RemediationAdvice string    `gorm:"type:text" json:"remediationAdvice"`
	FixedCodeSnippet  string    `gorm:"type:text" json:"fixedCodeSnippet"`
	Status            string    `gorm:"size:20;not null;index" json:"status"`
	CreatedAt         time.Time `json:"createdAt"`
}

func (f *Finding) BeforeCreate(_ *gorm.DB) error {
	if f.ID == "" {
		f.ID = uuid.NewString()
	}
	return nil
}

type CodeFile struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}
