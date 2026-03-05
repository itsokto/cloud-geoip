package writer

import (
	"os"
	"path/filepath"
	"strings"
)

var _ Writer = (*PlainWriter)(nil)

type PlainWriter struct{}

func (w *PlainWriter) Write(baseDir string, entries []Entry) error {
	dir, err := ensureDir(baseDir, "plain")
	if err != nil {
		return err
	}

	for _, e := range entries {
		path := filepath.Join(dir, e.Name+".txt")
		if err := os.WriteFile(path, []byte(strings.Join(e.Prefixes, "\n")+"\n"), 0644); err != nil {
			return err
		}
	}
	return nil
}
