package democorp

import "os"

func ReadUpload(name string) ([]byte, error) {
	return os.ReadFile("/var/app/uploads/" + name)
}
