package app

import (
	"app/clock"
	"app/tests/testutils"
	"reflect"
	"testing"
	"time"
)

type Mock struct {
	App
	FakeCreateBlock func(nonce int, previousHash string) *Block
}

func (m *Mock) CreateBlock(nonce int, previousHash string) *Block {
	return m.FakeCreateBlock(nonce, previousHash)
}

func TestNewBlock(t *testing.T) {
	t.Parallel()
	testutils.SetFakeTime(time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local))
	got := NewBlock(0, "init hash")
	want := &Block{
		timestamp:    clock.Now().UnixNano(),
		nonce:        0,
		previousHash: "init hash",
		transactions: []string{},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("テストに失敗しました。")
	}
}

func TestNewBlockchain(t *testing.T) {
	t.Parallel()
	testutils.SetFakeTime(time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local))
	chain := []*Block{{
		timestamp:    clock.Now().UnixNano(),
		nonce:        0,
		previousHash: "init hash",
		transactions: []string{},
	}}
	got := NewBlockChain()
	want := &Blockchain{
		transactionPool: []string{},
		chain:           chain,
	}
	if !reflect.DeepEqual(&got.chain, &want.chain) {
		t.Errorf("テストに失敗しました。")
	}
	if &got.transactionPool == &want.transactionPool {
		t.Errorf("テストに失敗しました。")
	}
}

func TestMain(t *testing.T) {
	t.Parallel()
	testutils.SetFakeTime(time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local))
	got := testutils.ExtractStdout(t, main)
	want := "nonce        :0\npreviousHash :init hash\ntimestamp    :1577804400000000000\ntransactions :[]"
	if got != want {
		t.Errorf("テストに失敗しました。 got: %s, want: %s", got, want)
	}
}
