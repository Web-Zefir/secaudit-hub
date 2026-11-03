package sast

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDemoPresetCoversAllRules(t *testing.T) {
	files, err := DemoFiles()
	if err != nil {
		t.Fatal(err)
	}
	findings := Scan(files)
	got := map[string]Finding{}
	for _, f := range findings {
		got[f.RuleID] = f
	}
	want := []string{
		"SAST-SEC-001", "SAST-SEC-002", "SAST-SQL-001", "SAST-CMD-001",
		"SAST-XSS-001", "SAST-EVAL-001", "SAST-PATH-001", "SAST-CRYP-001",
		"SAST-CRYP-002", "SAST-NET-001", "SAST-HTTP-001", "SAST-CORS-001",
	}
	for _, id := range want {
		f, ok := got[id]
		if !ok {
			t.Errorf("demo preset missed %s", id)
			continue
		}
		if f.LineNumber < 1 || f.CodeSnippet == "" || f.FixedCodeSnippet == "" || f.VulnerableLine == "" {
			t.Errorf("%s missing context or fix", id)
		}
		if f.Category == "" || f.Severity == "" {
			t.Errorf("%s missing classification", id)
		}
	}
}

func TestRulesDetectTargetedSamples(t *testing.T) {
	cases := []struct {
		name string
		file SourceFile
		rule string
	}{
		{"aws key", SourceFile{"cfg.env", "AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE\n"}, "SAST-SEC-001"},
		{"jwt", SourceFile{"a.env", "JWT_SECRET=prod-jwt-signing-secret-9f3a\n"}, "SAST-SEC-001"},
		{"private key", SourceFile{"k.pem", "-----BEGIN RSA PRIVATE KEY-----\n"}, "SAST-SEC-001"},
		{"dsn", SourceFile{"db.go", "dsn := \"postgres://user:secret@db:5432/app\"\n"}, "SAST-SEC-002"},
		{"sprintf sql", SourceFile{"q.go", "q := fmt.Sprintf(\"SELECT id FROM users WHERE email = '%s'\", email)\n"}, "SAST-SQL-001"},
		{"f string", SourceFile{"q.py", "cur.execute(f\"SELECT * FROM users WHERE id = {user_id}\")\n"}, "SAST-SQL-001"},
		{"shell", SourceFile{"c.go", "exec.Command(\"bash\", \"-c\", \"echo \"+name)\n"}, "SAST-CMD-001"},
		{"node exec", SourceFile{"a.js", "child_process.exec(userInput)\n"}, "SAST-CMD-001"},
		{"inner html", SourceFile{"a.tsx", "node.innerHTML = bio\n"}, "SAST-XSS-001"},
		{"eval", SourceFile{"a.py", "return eval(user_code)\n"}, "SAST-EVAL-001"},
		{"function ctor", SourceFile{"a.js", "const fn = new Function(code)\n"}, "SAST-EVAL-001"},
		{"path", SourceFile{"a.go", "os.ReadFile(\"/tmp/\" + name)\n"}, "SAST-PATH-001"},
		{"md5", SourceFile{"a.go", "sum := md5.Sum([]byte(password))\n"}, "SAST-CRYP-001"},
		{"sha1", SourceFile{"a.js", "crypto.createHash(\"sha1\").update(password)\n"}, "SAST-CRYP-001"},
		{"rand", SourceFile{"a.go", "import \"math/rand\"\nfunc NewSessionToken() string {\n  return fmt.Sprintf(\"%d\", rand.Intn(10))\n}\n"}, "SAST-CRYP-002"},
		{"js random", SourceFile{"a.js", "function sessionToken() {\n  return Math.random().toString()\n}\n"}, "SAST-CRYP-002"},
		{"tls", SourceFile{"a.go", "cfg := &tls.Config{InsecureSkipVerify: true}\n"}, "SAST-NET-001"},
		{"listen", SourceFile{"a.go", "http.ListenAndServe(\":8080\", mux)\n"}, "SAST-HTTP-001"},
		{"cors", SourceFile{"a.go", "w.Header().Set(\"Access-Control-Allow-Origin\", \"*\")\nw.Header().Set(\"Access-Control-Allow-Credentials\", \"true\")\n"}, "SAST-CORS-001"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if !hasRule(Scan([]SourceFile{tc.file}), tc.rule) {
				t.Fatalf("expected %s", tc.rule)
			}
		})
	}
}

func TestRulesIgnoreSafeSamples(t *testing.T) {
	files := []SourceFile{
		{"env.go", "apiKey := os.Getenv(\"API_KEY\")\n"},
		{"q.go", "db.Query(\"SELECT id FROM users WHERE email = $1\", email)\n"},
		{"git.go", "exec.Command(\"git\", \"status\")\n"},
		{"hash.go", "sum := sha256.Sum256([]byte(password))\n"},
		{"rand.go", "import \"crypto/rand\"\nfunc token() { rand.Read(buf) }\n"},
		{"srv.go", "srv := &http.Server{Addr: \":8080\", ReadHeaderTimeout: 5 * time.Second}\n"},
		{"path.go", "os.ReadFile(filepath.Clean(\"/tmp/\" + name))\n"},
	}
	findings := Scan(files)
	if len(findings) != 0 {
		t.Fatalf("safe samples produced findings: %+v", findings)
	}
}

func TestSummarizePenalties(t *testing.T) {
	sum := Summarize([]Finding{{Severity: "CRITICAL"}, {Severity: "HIGH"}, {Severity: "MEDIUM"}, {Severity: "LOW"}})
	if sum.Score != 56 || sum.Critical != 1 || sum.High != 1 || sum.Medium != 1 || sum.Low != 1 {
		t.Fatalf("unexpected summary: %+v", sum)
	}
	heavy := Summarize([]Finding{{Severity: "CRITICAL"}, {Severity: "CRITICAL"}, {Severity: "CRITICAL"}, {Severity: "CRITICAL"}, {Severity: "CRITICAL"}})
	if heavy.Score != 0 {
		t.Fatalf("score = %d", heavy.Score)
	}
	if Summarize(nil).Score != 100 {
		t.Fatal("empty score")
	}
}

func TestValidateGitURL(t *testing.T) {
	if err := ValidateGitURL("https://github.com/example/app.git"); err != nil {
		t.Fatal(err)
	}
	for _, raw := range []string{"http://github.com/a/b", "https://user:token@github.com/a/b", "file:///tmp/repo", "https://127.0.0.1/a"} {
		if err := ValidateGitURL(raw); err == nil {
			t.Fatalf("accepted %s", raw)
		}
	}
}

func TestCollectFilesSkipsVendor(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "vendor"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "main.go"), []byte("package main\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "vendor", "lib.go"), []byte("package vendor\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	files, err := CollectFiles(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 || files[0].Path != "main.go" {
		t.Fatalf("files = %+v", files)
	}
}

func hasRule(findings []Finding, id string) bool {
	for _, f := range findings {
		if f.RuleID == id {
			return true
		}
	}
	return false
}
