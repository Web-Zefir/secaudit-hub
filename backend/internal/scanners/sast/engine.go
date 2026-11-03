package sast

import (
	"fmt"
	"strings"
)

const Name = "SAST"

type SourceFile struct {
	Path    string
	Content string
}

type Finding struct {
	RuleID            string
	Category          string
	Severity          string
	Title             string
	Description       string
	FilePath          string
	LineNumber        int
	CodeSnippet       string
	VulnerableLine    string
	RemediationAdvice string
	FixedCodeSnippet  string
}

type Summary struct {
	Score    int
	Critical int
	High     int
	Medium   int
	Low      int
}

type hit struct {
	line int
}

type rule struct {
	id          string
	severity    string
	cwe         string
	owasp       string
	title       string
	description string
	remediation string
	match       func(file SourceFile, lines []string) []hit
}

// Scan runs every built-in rule against the provided sources.
func Scan(files []SourceFile) []Finding {
	rules := builtinRules()
	out := make([]Finding, 0)
	seen := make(map[string]struct{})
	for _, file := range files {
		lines := splitLines(file.Content)
		for _, rl := range rules {
			for _, h := range rl.match(file, lines) {
				if h.line < 1 || h.line > len(lines) {
					continue
				}
				key := rl.id + "|" + file.Path + "|" + fmt.Sprint(h.line)
				if _, ok := seen[key]; ok {
					continue
				}
				seen[key] = struct{}{}
				vuln := strings.TrimRight(lines[h.line-1], "\r")
				out = append(out, Finding{
					RuleID:            rl.id,
					Category:          rl.owasp + " · " + rl.cwe,
					Severity:          rl.severity,
					Title:             rl.title,
					Description:       rl.description,
					FilePath:          file.Path,
					LineNumber:        h.line,
					CodeSnippet:       snippet(lines, h.line),
					VulnerableLine:    vuln,
					RemediationAdvice: rl.remediation,
					FixedCodeSnippet:  fixedSnippet(rl.id),
				})
			}
		}
	}
	return out
}

// Summarize applies a fixed penalty per severity and clamps the score to 0..100.
// CRITICAL costs 25, HIGH 12, MEDIUM 5, and LOW 2. A clean tree stays at 100.
func Summarize(findings []Finding) Summary {
	var sum Summary
	for _, f := range findings {
		switch f.Severity {
		case "CRITICAL":
			sum.Critical++
		case "HIGH":
			sum.High++
		case "MEDIUM":
			sum.Medium++
		case "LOW":
			sum.Low++
		}
	}
	penalty := sum.Critical*25 + sum.High*12 + sum.Medium*5 + sum.Low*2
	if penalty > 100 {
		penalty = 100
	}
	sum.Score = 100 - penalty
	return sum
}

func splitLines(content string) []string {
	return strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")
}

func snippet(lines []string, line int) string {
	idx := line - 1
	start := idx - 2
	if start < 0 {
		start = 0
	}
	end := idx + 3
	if end > len(lines) {
		end = len(lines)
	}
	var b strings.Builder
	for i := start; i < end; i++ {
		marker := " "
		if i == idx {
			marker = ">"
		}
		fmt.Fprintf(&b, "%s %4d | %s\n", marker, i+1, strings.TrimRight(lines[i], "\r"))
	}
	return strings.TrimRight(b.String(), "\n")
}

func isComment(line string) bool {
	t := strings.TrimSpace(line)
	if strings.HasPrefix(t, "//") || strings.HasPrefix(t, "#") || strings.HasPrefix(t, "/*") || strings.HasPrefix(t, "*") {
		return true
	}
	return strings.HasPrefix(t, "--") && (len(t) == 2 || t[2] == ' ' || t[2] == '\t')
}
