package main

import (
	"os"
	"path/filepath"
	"strings"
)

func addSuffixIfMissing(s *string, suffix string) {
	if !strings.HasSuffix(*s, suffix) {
		*s = *s + suffix
	}
}

func replaceExt(path string, newExt string) string {
	return strings.TrimSuffix(path, filepath.Ext(path)) + newExt
}

func dirExists(path string) (bool, error) {
	dir, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer dir.Close()

	if stat, err := dir.Stat(); err != nil || !stat.IsDir() {
		return false, err
	}
	return true, nil
}

func IntMax(a int, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}
