package dast

// Name identifies the dynamic analysis module reserved for a later epic.
const Name = "DAST"

// Ready reports that DAST is not executable in this build.
func Ready() bool { return false }
