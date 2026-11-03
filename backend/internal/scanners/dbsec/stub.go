package dbsec

// Name identifies the database security module reserved for a later epic.
const Name = "DB Security"

// Ready reports that database security checks are not executable in this build.
func Ready() bool { return false }
