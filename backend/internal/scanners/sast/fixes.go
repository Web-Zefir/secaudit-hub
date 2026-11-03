package sast

func fixedSnippet(ruleID string) string {
	switch ruleID {
	case "SAST-SEC-001":
		return `secret := os.Getenv("APP_SECRET")
if secret == "" {
    return fmt.Errorf("APP_SECRET is required")
}`
	case "SAST-SEC-002":
		return `dsn := os.Getenv("DATABASE_URL")
if dsn == "" {
    return fmt.Errorf("DATABASE_URL is required")
}`
	case "SAST-SQL-001":
		return `rows, err := db.Query(
    "SELECT id, email, role FROM users WHERE email = $1",
    email,
)`
	case "SAST-CMD-001":
		return `cmd := exec.Command("pg_dump", "-h", host, "billing")
cmd.Env = os.Environ()
if err := cmd.Run(); err != nil {
    return err
}`
	case "SAST-XSS-001":
		return `return <p>{bio}</p>`
	case "SAST-EVAL-001":
		return `switch action {
case "ping":
    return ping(), nil
default:
    return nil, fmt.Errorf("unsupported action")
}`
	case "SAST-PATH-001":
		return `clean := filepath.Clean("/" + name)
full := filepath.Join(uploadRoot, clean)
rel, err := filepath.Rel(uploadRoot, full)
if err != nil || strings.HasPrefix(rel, "..") {
    return nil, fmt.Errorf("path escapes upload root")
}
return os.ReadFile(full)`
	case "SAST-CRYP-001":
		return `sum := sha256.Sum256([]byte(password))
stored := hex.EncodeToString(sum[:])
// Prefer bcrypt or argon2id for passwords.
hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)`
	case "SAST-CRYP-002":
		return `buf := make([]byte, 32)
if _, err := cryptorand.Read(buf); err != nil {
    return "", err
}
return hex.EncodeToString(buf), nil`
	case "SAST-NET-001":
		return `tlsConfig := &tls.Config{MinVersion: tls.VersionTLS12}`
	case "SAST-HTTP-001":
		return `srv := &http.Server{
    Addr:              ":8080",
    Handler:           mux,
    ReadHeaderTimeout: 5 * time.Second,
    ReadTimeout:       15 * time.Second,
    WriteTimeout:      15 * time.Second,
    IdleTimeout:       60 * time.Second,
}
return srv.ListenAndServe()`
	case "SAST-CORS-001":
		return `origin := r.Header.Get("Origin")
if origin == allowedOrigin {
    w.Header().Set("Access-Control-Allow-Origin", origin)
    w.Header().Set("Vary", "Origin")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
}`
	default:
		return ""
	}
}
