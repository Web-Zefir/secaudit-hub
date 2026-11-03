package iac

// Name identifies the infrastructure-as-code module reserved for a later epic.
const Name = "IaC"

// Ready reports that IaC scanning is not executable in this build.
func Ready() bool { return false }
