package sast

import (
	"regexp"
	"strings"
)

func builtinRules() []rule {
	return []rule{
		secretRule(),
		dbSecretRule(),
		sqlRule(),
		cmdRule(),
		xssRule(),
		evalRule(),
		pathRule(),
		weakHashRule(),
		prngRule(),
		tlsRule(),
		httpTimeoutRule(),
		corsRule(),
	}
}

func secretRule() rule {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`AKIA[0-9A-Z]{16}`),
		regexp.MustCompile(`(?i)aws_secret_access_key\s*[:=]\s*\S+`),
		regexp.MustCompile(`(?i)\b(jwt[_-]?secret|api[_-]?key|api[_-]?token|secret[_-]?key)\b\s*[:=]\s*(?:['\"][^'\"]{8,}['\"]|[A-Za-z0-9_./+=-]{12,})`),
		regexp.MustCompile(`sk_live_[0-9A-Za-z]{8,}`),
		regexp.MustCompile(`-----BEGIN (RSA |OPENSSH |EC |DSA )?PRIVATE KEY-----`),
	}
	return rule{
		id: "SAST-SEC-001", severity: "CRITICAL", cwe: "CWE-798", owasp: "A07:2021",
		title:       "Hardcoded credential or private key",
		description: "A cloud key, API token, JWT secret, or private key is committed in source. Anyone with repository access can reuse it.",
		remediation: "Load secrets from a secret manager or environment variables and rotate any exposed credential.",
		match:       matchLinePatterns(patterns),
	}
}

func dbSecretRule() rule {
	dsn := regexp.MustCompile(`(?i)(postgres|postgresql|mysql|mongodb(\+srv)?|redis|amqp)://[^:\s]+:[^@\s]+@`)
	named := regexp.MustCompile(`(?i)\b(db_password|database_password|mysql_password|postgres_password|redis_password)\b\s*[:=]\s*['\"]?[^\s'\"]{4,}`)
	assign := regexp.MustCompile(`(?i)\b(password|passwd|pwd)\b\s*[:=]\s*['\"][^'\"]{4,}['\"]`)
	dbctx := regexp.MustCompile(`(?i)postgres|mysql|mongodb|redis|database|dsn|connection`)
	return rule{
		id: "SAST-SEC-002", severity: "HIGH", cwe: "CWE-259", owasp: "A07:2021",
		title:       "Hardcoded database password or connection string",
		description: "A database credential is embedded in source. The password travels with the code and cannot be rotated independently.",
		remediation: "Build the DSN from environment variables and keep the password outside the repository.",
		match: func(file SourceFile, lines []string) []hit {
			var hits []hit
			for i, line := range lines {
				if isComment(line) {
					continue
				}
				if dsn.MatchString(line) || named.MatchString(line) || (assign.MatchString(line) && dbctx.MatchString(line)) {
					hits = append(hits, hit{line: i + 1})
				}
			}
			return hits
		},
	}
}

func sqlRule() rule {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)fmt\.Sprintf\(\s*(?:"[^"]*(?:SELECT|INSERT|UPDATE|DELETE)|'[^']*(?:SELECT|INSERT|UPDATE|DELETE)|` + "`[^`]*(?:SELECT|INSERT|UPDATE|DELETE))"),
		regexp.MustCompile(`(?i)(?:\"[^\"]*(?:SELECT|INSERT|UPDATE|DELETE)[^\"]*\"|'[^']*(?:SELECT|INSERT|UPDATE|DELETE)[^']*'|` + "`[^`]*(?:SELECT|INSERT|UPDATE|DELETE)[^`]*`" + `)\s*\+`),
		regexp.MustCompile(`(?i)f["'][^"']*(?:SELECT|INSERT|UPDATE|DELETE|WHERE)`),
		regexp.MustCompile(`(?i)["'][^"']*(?:SELECT|INSERT|UPDATE|DELETE)[^"']*["']\.format\(`),
		regexp.MustCompile(`(?i)\.execute\(\s*f?["'].*(?:SELECT|INSERT|UPDATE|DELETE)`),
	}
	return rule{
		id: "SAST-SQL-001", severity: "CRITICAL", cwe: "CWE-89", owasp: "A03:2021",
		title:       "SQL injection via dynamic query construction",
		description: "SQL is built with string concatenation, formatting, or an f-string. Untrusted input can change the query structure.",
		remediation: "Use prepared statements with bound parameters. Never concatenate user input into SQL text.",
		match:       matchLinePatterns(patterns),
	}
}

func cmdRule() rule {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`exec\.Command\(\s*"(sh|bash|cmd|powershell)"\s*,\s*"-(c|C)"`),
		regexp.MustCompile(`exec\.Command\([^)]*\+`),
		regexp.MustCompile(`child_process\.exec(Sync)?\(`),
		regexp.MustCompile(`\bos\.system\(`),
		regexp.MustCompile(`subprocess\.(run|call|Popen|check_output|check_call)\([^)]*shell\s*=\s*True`),
	}
	return rule{
		id: "SAST-CMD-001", severity: "CRITICAL", cwe: "CWE-78", owasp: "A03:2021",
		title:       "OS command injection",
		description: "A shell is invoked with a command string that can include untrusted data. Metacharacters can execute extra commands.",
		remediation: "Call the binary directly with a fixed argument vector. Do not pass user input through a shell.",
		match:       matchLinePatterns(patterns),
	}
}

func xssRule() rule {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`dangerouslySetInnerHTML`),
		regexp.MustCompile(`\.innerHTML\s*=`),
		regexp.MustCompile(`document\.write\(`),
	}
	return rule{
		id: "SAST-XSS-001", severity: "HIGH", cwe: "CWE-79", owasp: "A03:2021",
		title:       "Unescaped HTML injection",
		description: "HTML is written into the DOM without escaping. Script or markup in the data will execute in the user's browser.",
		remediation: "Render text through the framework's default escaping. Sanitize HTML with a maintained allow-list if markup is required.",
		match:       matchLinePatterns(patterns),
	}
}

func evalRule() rule {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`\beval\s*\(`),
		regexp.MustCompile(`\bnew\s+Function\s*\(`),
	}
	return rule{
		id: "SAST-EVAL-001", severity: "HIGH", cwe: "CWE-95", owasp: "A03:2021",
		title:       "Dynamic code execution",
		description: "eval or the Function constructor executes a string as code. If any part of that string is user-controlled, it is remote code execution.",
		remediation: "Replace dynamic evaluation with an explicit allow-list of operations or a restricted interpreter.",
		match:       matchLinePatterns(patterns),
	}
}

func matchLinePatterns(patterns []*regexp.Regexp) func(SourceFile, []string) []hit {
	return func(_ SourceFile, lines []string) []hit {
		var hits []hit
		for i, line := range lines {
			if isComment(line) {
				continue
			}
			for _, p := range patterns {
				if p.MatchString(line) {
					hits = append(hits, hit{line: i + 1})
					break
				}
			}
		}
		return hits
	}
}

func pathRule() rule {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`os\.(ReadFile|Open|WriteFile|Create|Remove|RemoveAll)\([^)]*\+`),
		regexp.MustCompile(`ioutil\.(ReadFile|WriteFile)\([^)]*\+`),
		regexp.MustCompile(`fs\.(readFile|readFileSync|writeFile|writeFileSync)\([^)]*\+`),
		regexp.MustCompile(`\bopen\([^)]*\+`),
	}
	return rule{
		id: "SAST-PATH-001", severity: "HIGH", cwe: "CWE-22", owasp: "A01:2021",
		title:       "Path traversal in file access",
		description: "A filesystem path is built by concatenation and is not passed through filepath.Clean or an equivalent root check.",
		remediation: "Clean the path, join it under a fixed root, and reject results that escape that root.",
		match: func(_ SourceFile, lines []string) []hit {
			var hits []hit
			for i, line := range lines {
				if isComment(line) || strings.Contains(line, "filepath.Clean") || strings.Contains(line, "path.Clean") {
					continue
				}
				for _, p := range patterns {
					if p.MatchString(line) {
						hits = append(hits, hit{line: i + 1})
						break
					}
				}
			}
			return hits
		},
	}
}
