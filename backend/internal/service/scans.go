package service

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/Web-Zefir/secaudit-hub/backend/internal/models"
	"github.com/Web-Zefir/secaudit-hub/backend/internal/repository"
	"github.com/Web-Zefir/secaudit-hub/backend/internal/scanners/sast"
	"github.com/Web-Zefir/secaudit-hub/backend/internal/worker"
)

type ScanService struct {
	projects *repository.ProjectRepo
	scans    *repository.ScanRepo
	findings *repository.FindingRepo
	pool     *worker.Pool
}

func NewScanService(projects *repository.ProjectRepo, scans *repository.ScanRepo, findings *repository.FindingRepo) *ScanService {
	return &ScanService{projects: projects, scans: scans, findings: findings, pool: worker.NewPool(3)}
}

type StartScanInput struct {
	TargetSource string
	RepoURL      string
	Files        []models.CodeFile
}

func (s *ScanService) Start(ctx context.Context, userID, projectID string, in StartScanInput) (*models.ScanJob, error) {
	if _, err := s.projects.ByID(ctx, userID, projectID); err != nil {
		return nil, err
	}
	source := strings.TrimSpace(in.TargetSource)
	ref, err := normalizeTarget(source, in)
	if err != nil {
		return nil, err
	}
	job := &models.ScanJob{
		ProjectID:    projectID,
		ScannerType:  models.ScannerSAST,
		TargetSource: source,
		TargetRef:    ref,
		Status:       models.ScanPending,
	}
	if err := s.scans.Create(ctx, job); err != nil {
		return nil, err
	}
	snapshot := *job
	files := append([]models.CodeFile(nil), in.Files...)
	repoURL := strings.TrimSpace(in.RepoURL)
	s.pool.Go(func() { s.execute(snapshot, files, repoURL) })
	return job, nil
}

func (s *ScanService) List(ctx context.Context, userID, projectID string) ([]models.ScanJob, error) {
	if _, err := s.projects.ByID(ctx, userID, projectID); err != nil {
		return nil, err
	}
	items, err := s.scans.ListByProject(ctx, userID, projectID)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []models.ScanJob{}
	}
	return items, nil
}

func (s *ScanService) Get(ctx context.Context, userID, scanID string) (*models.ScanJob, []models.Finding, error) {
	job, err := s.scans.ByID(ctx, userID, scanID)
	if err != nil {
		return nil, nil, err
	}
	findings, err := s.findings.ByScan(ctx, userID, scanID)
	if err != nil {
		return nil, nil, err
	}
	if findings == nil {
		findings = []models.Finding{}
	}
	return job, findings, nil
}

func (s *ScanService) UpdateFindingStatus(ctx context.Context, userID, findingID, status string) (*models.Finding, error) {
	switch status {
	case models.FindingOpen, models.FindingResolved, models.FindingFalsePositive:
	default:
		return nil, ErrInvalidInput
	}
	return s.findings.UpdateStatus(ctx, userID, findingID, status)
}

func normalizeTarget(source string, in StartScanInput) (string, error) {
	switch source {
	case models.TargetDemo:
		return "Vulnerable demo corpus", nil
	case models.TargetCustom:
		if len(in.Files) == 0 || len(in.Files) > 20 {
			return "", ErrInvalidInput
		}
		names := make([]string, 0, len(in.Files))
		for _, file := range in.Files {
			if _, err := cleanUploadName(file.Name); err != nil {
				return "", err
			}
			if strings.TrimSpace(file.Content) == "" || len(file.Content) > 200_000 {
				return "", ErrInvalidInput
			}
			names = append(names, file.Name)
		}
		return strings.Join(names, ", "), nil
	case models.TargetGit:
		if err := sast.ValidateGitURL(in.RepoURL); err != nil {
			return "", fmt.Errorf("%w: %s", ErrInvalidInput, err.Error())
		}
		return strings.TrimSpace(in.RepoURL), nil
	default:
		return "", ErrInvalidInput
	}
}

func cleanUploadName(name string) (string, error) {
	name = strings.TrimSpace(strings.ReplaceAll(name, "\\", "/"))
	if name == "" || strings.Contains(name, "..") || path.IsAbs(name) {
		return "", ErrInvalidInput
	}
	clean := path.Clean(name)
	if clean == "." || strings.HasPrefix(clean, "../") || clean == ".." {
		return "", ErrInvalidInput
	}
	return clean, nil
}

func (s *ScanService) execute(job models.ScanJob, files []models.CodeFile, repoURL string) {
	ctx := context.Background()
	started := time.Now()
	if err := s.scans.MarkRunning(ctx, job.ID, started); err != nil {
		return
	}
	defer func() {
		if recover() != nil {
			s.fail(ctx, &job, started, "scan worker panicked")
		}
	}()
	sources, cleanup, err := resolveSources(ctx, job.TargetSource, files, repoURL)
	if cleanup != nil {
		defer cleanup()
	}
	if err != nil {
		s.fail(ctx, &job, started, err.Error())
		return
	}
	if len(sources) == 0 {
		s.fail(ctx, &job, started, "no supported source files were found")
		return
	}
	raw := sast.Scan(sources)
	summary := sast.Summarize(raw)
	stored := make([]models.Finding, 0, len(raw))
	for _, f := range raw {
		stored = append(stored, models.Finding{
			ScanJobID: job.ID, ProjectID: job.ProjectID, RuleID: f.RuleID, Category: f.Category,
			Severity: f.Severity, Title: f.Title, Description: f.Description, FilePath: f.FilePath,
			LineNumber: f.LineNumber, CodeSnippet: f.CodeSnippet, VulnerableLine: f.VulnerableLine,
			RemediationAdvice: f.RemediationAdvice, FixedCodeSnippet: f.FixedCodeSnippet,
			Status: models.FindingOpen,
		})
	}
	if err := s.findings.CreateBatch(ctx, stored); err != nil {
		s.fail(ctx, &job, started, err.Error())
		return
	}
	finished := time.Now()
	job.Status = models.ScanCompleted
	job.SecurityScore = summary.Score
	job.FilesScanned = len(sources)
	job.DurationMs = finished.Sub(started).Milliseconds()
	job.CriticalCount = summary.Critical
	job.HighCount = summary.High
	job.MediumCount = summary.Medium
	job.LowCount = summary.Low
	job.FinishedAt = &finished
	if err := s.scans.Finish(ctx, &job); err != nil {
		return
	}
	_ = s.projects.TouchScan(ctx, job.ProjectID, summary.Score, finished)
}

func (s *ScanService) fail(ctx context.Context, job *models.ScanJob, started time.Time, message string) {
	finished := time.Now()
	job.Status = models.ScanFailed
	job.ErrorMessage = message
	job.DurationMs = finished.Sub(started).Milliseconds()
	job.FinishedAt = &finished
	_ = s.scans.Finish(ctx, job)
}

func resolveSources(ctx context.Context, source string, files []models.CodeFile, repoURL string) ([]sast.SourceFile, func(), error) {
	switch source {
	case models.TargetDemo:
		demo, err := sast.DemoFiles()
		return demo, nil, err
	case models.TargetCustom:
		out := make([]sast.SourceFile, 0, len(files))
		for _, file := range files {
			name, err := cleanUploadName(file.Name)
			if err != nil {
				return nil, nil, err
			}
			out = append(out, sast.SourceFile{Path: name, Content: file.Content})
		}
		return out, nil, nil
	case models.TargetGit:
		return sast.Clone(ctx, repoURL)
	default:
		return nil, nil, errors.New("unsupported target")
	}
}
