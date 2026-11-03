package modules

import (
	"github.com/Web-Zefir/secaudit-hub/backend/internal/scanners/dast"
	"github.com/Web-Zefir/secaudit-hub/backend/internal/scanners/dbsec"
	"github.com/Web-Zefir/secaudit-hub/backend/internal/scanners/iac"
	"github.com/Web-Zefir/secaudit-hub/backend/internal/scanners/sast"
)

type Info struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

func All() []Info {
	return []Info{
		{ID: "sast", Name: sast.Name, Status: status(true)},
		{ID: "dast", Name: dast.Name, Status: status(dast.Ready())},
		{ID: "iac", Name: iac.Name, Status: status(iac.Ready())},
		{ID: "dbsec", Name: dbsec.Name, Status: status(dbsec.Ready())},
	}
}

func status(ready bool) string {
	if ready {
		return "available"
	}
	return "roadmap"
}
