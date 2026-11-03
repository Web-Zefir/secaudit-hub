package sast

import (
	"regexp"
	"strings"
)

func weakHashRule() rule {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`\bmd5\.(New|Sum)\(`),
		regexp.MustCompile(`\bsha1\.(New|Sum)\(`),
		regexp.MustCompile(`hashlib\.(md5|sha1)\(`),
		regexp.MustCompile(`createHash\(\s*['\"](md5|sha1)['\"]`),
	}
	return rule{
		id: "SAST-CRYP-001", severity: "MEDIUM", cwe: "CWE-327", owasp: "A02:2021",
		title:       "Weak hash algorithm for sensitive data",
		description: "MD5 or SHA-1 is used. Both are cheap to collide and are not suitable for passwords, tokens, or integrity of security decisions.",
		remediation: "Use bcrypt or argon2id for passwords and SHA-256 or SHA-512 for general integrity checks.",
		match:       matchLinePatterns(patterns),
	}
}

func prngRule() rule {
	mathUse := regexp.MustCompile(`\brand\.(Intn|Int|Int31|Float64|Read|Perm|ExpFloat64)\(`)
	jsToken := regexp.MustCompile(`(?i)token|secret|password|session|nonce|apiKey`)
	return rule{
		id: "SAST-CRYP-002", severity: "MEDIUM", cwe: "CWE-338", owasp: "A02:2021",
		title:       "Insecure random generator for a security value",
		description: "math/rand or Math.random is not a cryptographic generator. Its output can be predicted and must not mint tokens or secrets.",
		remediation: "Use crypto/rand or Web Crypto getRandomValues for tokens, session identifiers, and keys.",
		match: func(file SourceFile, lines []string) []hit {
			var hits []hit
			usesMathRand := strings.Contains(file.Content, `"math/rand"`)
			for i, line := range lines {
				if isComment(line) {
					continue
				}
				if usesMathRand && mathUse.MatchString(line) && windowHas(lines, i, 5, jsToken) {
					hits = append(hits, hit{line: i + 1})
					continue
				}
				if strings.Contains(line, "Math.random()") && windowHas(lines, i, 4, jsToken) {
					hits = append(hits, hit{line: i + 1})
				}
			}
			return hits
		},
	}
}

func tlsRule() rule {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`InsecureSkipVerify\s*:\s*true`),
		regexp.MustCompile(`rejectUnauthorized\s*:\s*false`),
		regexp.MustCompile(`ssl\._create_unverified_context\(`),
		regexp.MustCompile(`CERT_NONE`),
	}
	return rule{
		id: "SAST-NET-001", severity: "MEDIUM", cwe: "CWE-295", owasp: "A02:2021",
		title:       "TLS certificate verification disabled",
		description: "The client accepts any certificate. An on-path attacker can impersonate the remote service.",
		remediation: "Keep certificate verification enabled and pin or trust a private CA explicitly when needed.",
		match:       matchLinePatterns(patterns),
	}
}

func httpTimeoutRule() rule {
	listen := regexp.MustCompile(`http\.ListenAndServe\(`)
	server := regexp.MustCompile(`http\.Server\{`)
	return rule{
		id: "SAST-HTTP-001", severity: "LOW", cwe: "CWE-400", owasp: "A05:2021",
		title:       "HTTP server without ReadHeaderTimeout",
		description: "A server that does not limit header read time can hold connections open. Slowloris-style clients exhaust workers.",
		remediation: "Construct http.Server with ReadHeaderTimeout and related timeouts instead of the bare ListenAndServe helper.",
		match: func(file SourceFile, lines []string) []hit {
			if strings.Contains(file.Content, "ReadHeaderTimeout") {
				return nil
			}
			var hits []hit
			for i, line := range lines {
				if isComment(line) {
					continue
				}
				if listen.MatchString(line) || server.MatchString(line) {
					hits = append(hits, hit{line: i + 1})
				}
			}
			return hits
		},
	}
}

func corsRule() rule {
	cred := regexp.MustCompile(`(?i)Access-Control-Allow-Credentials["']?\s*,\s*["']true["']|Allow-Credentials["']?\s*[:=]\s*["']?true|AllowCredentials\s*:\s*true`)
	origin := regexp.MustCompile(`Access-Control-Allow-Origin["']?\s*,\s*["']\*["']|Access-Control-Allow-Origin\s*[:=]\s*["']?\*|AllowOrigins\s*:\s*\[\]string\{\s*"\*"\s*\}`)
	return rule{
		id: "SAST-CORS-001", severity: "LOW", cwe: "CWE-942", owasp: "A05:2021",
		title:       "Overly permissive CORS with credentials",
		description: "Any origin is allowed while credentials are enabled. A hostile site can make authenticated cross-origin requests.",
		remediation: "Reflect only an explicit allow-list of origins when credentials are required. Never combine * with credentials.",
		match: func(file SourceFile, lines []string) []hit {
			if !cred.MatchString(file.Content) {
				return nil
			}
			var hits []hit
			for i, line := range lines {
				if isComment(line) {
					continue
				}
				if origin.MatchString(line) {
					hits = append(hits, hit{line: i + 1})
				}
			}
			return hits
		},
	}
}

func windowHas(lines []string, idx, radius int, re *regexp.Regexp) bool {
	start := idx - radius
	if start < 0 {
		start = 0
	}
	end := idx + radius + 1
	if end > len(lines) {
		end = len(lines)
	}
	return re.MatchString(strings.Join(lines[start:end], "\n"))
}
