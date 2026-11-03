package sast

import (
	"embed"
	"io/fs"
	"sort"
	"strings"
)

//go:embed demodata/*
var demoFS embed.FS

func DemoFiles() ([]SourceFile, error) {
	entries, err := demoFS.ReadDir("demodata")
	if err != nil {
		return nil, err
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })
	files := make([]SourceFile, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		body, err := fs.ReadFile(demoFS, "demodata/"+entry.Name())
		if err != nil {
			return nil, err
		}
		files = append(files, SourceFile{Path: "demo/" + entry.Name(), Content: strings.ReplaceAll(string(body), "\r\n", "\n")})
	}
	return files, nil
}
