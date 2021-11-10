package main

import (
	"path/filepath"
	"testing"
)

func TestFileGlob(t *testing.T) {
	files, err := filepath.Glob("./test/*.go")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(files)
}
