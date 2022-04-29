package utils

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func ExtractStdout(t *testing.T, fnc func()) string {
	t.Helper()
	orgStdout := os.Stdout

	defer func() {
		os.Stdout = orgStdout
	}()

	r, w, _ := os.Pipe()
	os.Stdout = w

	fnc()

	w.Close()
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("Failed to read buf: %v", err)
	}
	return strings.TrimRight(buf.String(), "\n")
}
