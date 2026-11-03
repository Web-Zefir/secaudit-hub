package democorp

import (
	"database/sql"
	"fmt"
)

func FindUser(db *sql.DB, email string) (*sql.Rows, error) {
	q := fmt.Sprintf("SELECT id, email, role FROM users WHERE email = '%s'", email)
	return db.Query(q)
}
