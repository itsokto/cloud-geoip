package writer

import (
	"os"
	"path/filepath"
)

type Entry struct {
	Name     string
	Prefixes []string
}

// Writer writes IP prefix data in a specific format.
type Writer interface {
	Write(baseDir string, entries []Entry) error
}

func ensureDir(baseDir string, parts ...string) (string, error) {
	dir := filepath.Join(append([]string{baseDir}, parts...)...)
	return dir, os.MkdirAll(dir, 0755)
}
