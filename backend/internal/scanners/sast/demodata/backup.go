package democorp

import "os/exec"

func RunBackup(host string) error {
	cmd := exec.Command("sh", "-c", "pg_dump -h "+host+" billing")
	return cmd.Run()
}
