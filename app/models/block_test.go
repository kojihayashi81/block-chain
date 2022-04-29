package models

import (
	"app/clock"
	"app/tests/utils"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func setup() {
	utils.SetFakeTime(time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local))
}
func TestNewBlock(t *testing.T) {
	t.Parallel()
	setup()

	got := NewBlock(0, [32]byte{})
	want := &Block{
		timestamp:    clock.Now().UnixNano(),
		nonce:        0,
		previousHash: [32]byte{},
		transactions: []string{},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("テストに失敗しました。 got:%x want:%x", got, want)
	}
}

func TestHash(t *testing.T) {
	t.Parallel()
	setup()
	b := &Block{
		timestamp:    clock.Now().UnixNano(),
		nonce:        0,
		previousHash: [32]byte{},
		transactions: []string{},
	}
	j, _ := json.Marshal(b)
	got := b.Hash()
	want := sha256.Sum256([]byte(j))
	if got != want {
		t.Errorf("テストに失敗しました。 got:%s want:%s", got, want)
	}
}

func TestPrint(t *testing.T) {
	t.Parallel()
	setup()
	time := clock.Now().UnixNano()
	hash := [32]byte{}
	tr := []string{}
	b := &Block{
		timestamp:    time,
		nonce:        0,
		previousHash: hash,
		transactions: tr,
	}
	got := utils.ExtractStdout(t, b.Print)
	var want string
	want += fmt.Sprintf("nonce        :%d\n", 0)
	want += fmt.Sprintf("previousHash :%x\n", hash)
	want += fmt.Sprintf("timestamp    :%d\n", time)
	want += fmt.Sprintf("transactions :%s", tr)
	if got != want {
		t.Errorf("テストに失敗しました。 got:%s want:%s", got, want)
	}
}
