package app

import (
	"app/clock"
	"app/tests/testUtils"
	"bytes"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestNewBlock(t *testing.T) {
	t.Parallel()
	testUtils.SetFakeTime(time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local))
	got := NewBlock(0, "init hash")
	want := &Block{
		timestamp:    clock.Now().UnixNano(),
		nonce:        0,
		previousHash: "init hash",
		transactions: []string{},
	}
	if !reflect.DeepEqual(&got, &want) {
		t.Errorf("failed to test.")
	}
}

func TestNewBlockchain(t *testing.T) {
	t.Parallel()
	testUtils.SetFakeTime(time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local))
	got := NewBlockChain()
	want := "aa"
	// want := &Blockchain{
	// 	transactionPool: []string{},
	// 	chain:           &Block{},
	// }
	if !reflect.DeepEqual(&got, &want) {
		t.Errorf("failed to test.")
	}
}
func TestMain(t *testing.T) {
	t.Parallel()
	testUtils.SetFakeTime(time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local))
	got := extractStdout(t, main)
	want := "nonce        :0\npreviousHash :init hash\ntimestamp    :1577836800000000000\ntransactions :[]"
	if got != want {
		t.Errorf("failed to test. got: %s, want: %s", got, want)
	}
}

func extractStdout(t *testing.T, fnc func()) string {
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
