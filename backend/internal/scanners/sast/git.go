package sast

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	maxGitFiles = 400
	maxFileSize = 256 * 1024
)

var errWalkStop = errors.New("walk stop")

var allowedExt = map[string]struct{}{
	".go": {}, ".ts": {}, ".tsx": {}, ".js": {}, ".jsx": {}, ".mjs": {}, ".cjs": {},
	".py": {}, ".java": {}, ".rb": {}, ".php": {}, ".cs": {}, ".c": {}, ".h": {},
	".cpp": {}, ".cc": {}, ".rs": {}, ".kt": {}, ".swift": {}, ".sql": {},
	".yml": {}, ".yaml": {}, ".json": {}, ".env": {}, ".sh": {}, ".tf": {},
}

// ValidateGitURL accepts only credential-free https URLs.
func ValidateGitURL(raw string) error {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || u.Scheme != "https" || u.Host == "" || u.User != nil {
		return errors.New("only public https git URLs without credentials are allowed")
	}
	host := strings.ToLower(u.Hostname())
	switch host {
	case "localhost", "127.0.0.1", "::1", "0.0.0.0", "169.254.169.254", "metadata.google.internal":
		return errors.New("refusing local or metadata git host")
	}
	if ip := net.ParseIP(host); ip != nil && !isPublicIP(ip) {
		return errors.New("refusing non-public git host")
	}
	return nil
}

// Clone fetches a public repository and returns its source files.
func Clone(parent context.Context, raw string) ([]SourceFile, func(), error) {
	if err := ValidateGitURL(raw); err != nil {
		return nil, nil, err
	}
	if err := assertPublicDNS(raw); err != nil {
		return nil, nil, err
	}
	dir, err := os.MkdirTemp("", "secaudit-git-*")
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	ctx, cancel := context.WithTimeout(parent, 45*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", "-c", "protocol.file.allow=never", "clone", "--depth", "1", raw, dir)
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	if out, err := cmd.CombinedOutput(); err != nil {
		cleanup()
		msg := strings.TrimSpace(string(out))
		if msg == "" {
			msg = err.Error()
		}
		return nil, nil, fmt.Errorf("git clone failed: %s", msg)
	}
	files, err := CollectFiles(dir)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	return files, cleanup, nil
}

func CollectFiles(root string) ([]SourceFile, error) {
	files := make([]SourceFile, 0)
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		name := d.Name()
		if d.IsDir() {
			if name == ".git" || name == "node_modules" || name == "vendor" || name == "dist" || name == "build" || name == ".venv" || name == "venv" {
				return filepath.SkipDir
			}
			return nil
		}
		if d.Type()&os.ModeSymlink != 0 {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(name))
		if name != ".env" {
			if _, ok := allowedExt[ext]; !ok {
				return nil
			}
		}
		info, err := d.Info()
		if err != nil || info.Size() > maxFileSize || info.Size() == 0 {
			return nil
		}
		body, err := os.ReadFile(path)
		if err != nil || bytes.Contains(body, []byte{0}) {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return nil
		}
		files = append(files, SourceFile{Path: filepath.ToSlash(rel), Content: string(body)})
		if len(files) >= maxGitFiles {
			return errWalkStop
		}
		return nil
	})
	if err != nil && !errors.Is(err, errWalkStop) {
		return nil, err
	}
	return files, nil
}

func assertPublicDNS(raw string) error {
	u, err := url.Parse(raw)
	if err != nil {
		return err
	}
	if net.ParseIP(u.Hostname()) != nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	ips, err := net.DefaultResolver.LookupIPAddr(ctx, u.Hostname())
	if err != nil {
		return fmt.Errorf("resolve git host: %w", err)
	}
	for _, ip := range ips {
		if !isPublicIP(ip.IP) {
			return errors.New("git host resolves to a non-public address")
		}
	}
	return nil
}

func isPublicIP(ip net.IP) bool {
	return ip != nil && !ip.IsLoopback() && !ip.IsPrivate() && !ip.IsLinkLocalUnicast() &&
		!ip.IsLinkLocalMulticast() && !ip.IsUnspecified() && !ip.IsMulticast()
}
