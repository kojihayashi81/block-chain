package models

import (
	"app/clock"
	"app/tests/utils"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewBlockchain(t *testing.T) {
	t.Parallel()
	setup()
	b := &Block{
		timestamp:    clock.Now().UnixNano(),
		nonce:        0,
		previousHash: [32]byte{},
		transactions: []string{},
	}
	appendBlock := &Block{
		timestamp:    clock.Now().UnixNano(),
		nonce:        0,
		previousHash: b.Hash(),
		transactions: []string{},
	}
	c := []*Block{appendBlock}
	got := NewBlockChain()
	want := &Blockchain{
		transactionPool: []string{},
		chain:           c,
	}
	if !reflect.DeepEqual(got.transactionPool, want.transactionPool) {
		t.Errorf("テストに失敗しました。\n got: %s\nwant: %s", &got.transactionPool, &want.transactionPool)
	}
	opts := []cmp.Option{
		cmp.AllowUnexported(Block{}),
		cmpopts.IgnoreFields(Block{}, "previousHash"),
	}
	if diff := cmp.Diff(*got.chain[0], *appendBlock, opts...); diff != "" {
		t.Errorf("テストに失敗しました。\n (-got +want):\n%s", diff)
	}
}

func TestCreateBlock(t *testing.T) {
	t.Parallel()
	setup()
	bc := new(Blockchain)
	nonce := 0
	previousHash := [32]byte{}
	transactions := []string{}
	got := bc.CreateBlock(nonce, previousHash)
	want := &Block{
		timestamp:    clock.Now().UnixNano(),
		nonce:        nonce,
		previousHash: previousHash,
		transactions: transactions,
	}
	if !reflect.DeepEqual(&got, &want) {
		t.Errorf("テストに失敗しました。\n got: %#v\nwant: %#v", &got, &want)
	}
	if contain(bc.chain, want) {
		t.Errorf("テストに失敗しました。\n got: %#v\nwant: %#v", bc.chain, &want)
	}
}

func TestLastBlock(t *testing.T) {
	t.Parallel()
	setup()
	b1 := &Block{
		timestamp:    clock.Now().UnixNano(),
		nonce:        1,
		previousHash: [32]byte{},
		transactions: []string{},
	}
	b2 := &Block{
		timestamp:    clock.Now().UnixNano() + 1,
		nonce:        2,
		previousHash: b1.Hash(),
		transactions: []string{},
	}
	b3 := &Block{
		timestamp:    clock.Now().UnixNano() + 2,
		nonce:        3,
		previousHash: b2.Hash(),
		transactions: []string{},
	}

	for _, test := range []struct {
		title       string
		input_tp    []string
		input_chain []*Block
		output      *Block
	}{
		{"b1変数のブロックが戻る", []string{}, []*Block{b1}, b1},
		{"b2変数のブロックが戻る", []string{}, []*Block{b1, b2}, b2},
		{"b3変数のブロックが戻る", []string{}, []*Block{b1, b2, b3}, b3},
		{"b2変数のブロックが戻る", []string{}, []*Block{b3, b1, b2}, b2},
	} {
		t.Run("TestLastBlock: "+test.title, func(t *testing.T) {
			bc := &Blockchain{
				transactionPool: test.input_tp,
				chain:           test.input_chain,
			}
			got := bc.LastBlock()
			want := test.output
			opts := []cmp.Option{
				cmp.AllowUnexported(Block{}),
				cmpopts.IgnoreFields(Block{}, "previousHash"),
			}
			if diff := cmp.Diff(got, want, opts...); diff != "" {
				t.Errorf("テストに失敗しました。\n (-got +want):\n%s", diff)
			}
		})
	}
}

func TestBlockchanPrint(t *testing.T) {
	t.Parallel()
	setup()
	b1 := &Block{
		timestamp:    clock.Now().UnixNano(),
		nonce:        1,
		previousHash: [32]byte{},
		transactions: []string{},
	}
	b2 := &Block{
		timestamp:    clock.Now().UnixNano() + 1,
		nonce:        2,
		previousHash: b1.Hash(),
		transactions: []string{},
	}
	var want string

	for _, test := range []struct {
		title       string
		input_chain []*Block
		output      string
	}{
		{"ターミナルにチェーンの内容が出力される", []*Block{b1, b2}, want},
	} {
		t.Run("TestPrint: "+test.title, func(t *testing.T) {
			bc := &Blockchain{
				transactionPool: []string{},
				chain:           test.input_chain,
			}
			for i, b := range bc.chain {
				want += fmt.Sprintf("%s Chain %d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
				want += fmt.Sprintf("nonce        :%d\n", b.nonce)
				want += fmt.Sprintf("previousHash :%x\n", b.previousHash)
				want += fmt.Sprintf("timestamp    :%d\n", b.timestamp)
				want += fmt.Sprintf("transactions :%s\n", b.transactions)
				if i == len(bc.chain)-1 {
					want += fmt.Sprintf("%s", strings.Repeat("*", 50))
				} else {
					want += fmt.Sprintf("%s\n", strings.Repeat("*", 50))
				}
				fmt.Println()
			}
			got := utils.ExtractStdout(t, bc.Print)
			opts := []cmp.Option{
				cmp.AllowUnexported(Block{}),
				cmpopts.IgnoreFields(Block{}, "previousHash"),
			}
			if diff := cmp.Diff(got, want, opts...); diff != "" {
				t.Errorf("テストに失敗しました。\n (-got +want):\n%s", diff)
			}
		})
	}
}

func contain(slice []*Block, target interface{}) bool {
	for _, item := range slice {
		if item == target {
			return true
		}
	}
	return false
}
