package models

import (
	"app/clock"
	"app/tests/utils"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/big"
	mathrand "math/rand"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func setup() {
	utils.SetFakeTime(time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local))
}
func TestNewBlock(t *testing.T) {
	t.Parallel()
	setup()

	got := NewBlock(0, [32]byte{}, []*Transaction{})
	want := &Block{
		timestamp:    clock.Now().UnixNano(),
		nonce:        0,
		previousHash: [32]byte{},
		transactions: []*Transaction{},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("テストに失敗しました。 got:%#v want:%#v", got, want)
	}
}

func TestHash(t *testing.T) {
	t.Parallel()
	setup()
	b := &Block{
		timestamp:    clock.Now().UnixNano(),
		nonce:        0,
		previousHash: [32]byte{},
		transactions: []*Transaction{},
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
	t1 := &Transaction{
		senderBlockchainAddress:    "senderA",
		recipientBlockchainAddress: "recipientB",
		value:                      1.0,
	}
	t2 := &Transaction{
		senderBlockchainAddress:    "senderB",
		recipientBlockchainAddress: "recipientC",
		value:                      1.0,
	}
	t3 := &Transaction{
		senderBlockchainAddress:    "senderA",
		recipientBlockchainAddress: "recipientC",
		value:                      1.0,
	}
	time := clock.Now().UnixNano()
	hash := [32]byte{}
	tr := []*Transaction{t1, t2, t3}
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
	for i, t := range tr {
		want += fmt.Sprintf("%s\n", strings.Repeat("-", 40))
		want += fmt.Sprintf(" sender_blockchain_address      %s\n", t.senderBlockchainAddress)
		want += fmt.Sprintf(" recipient_blockchain_address   %s\n", t.recipientBlockchainAddress)
		if i == len(b.transactions)-1 {
			want += fmt.Sprintf(" value                          %.1f", t.value)
		} else {
			want += fmt.Sprintf(" value                          %.1f\n", t.value)
		}
	}
	if got != want {
		t.Errorf("テストに失敗しました。 got:%s want:%s", got, want)
	}
}

func TestNewBlockchain(t *testing.T) {
	t.Parallel()
	setup()
	b := &Block{
		timestamp:    clock.Now().UnixNano(),
		nonce:        0,
		previousHash: [32]byte{},
		transactions: []*Transaction{},
	}
	appendBlock := &Block{
		timestamp:    clock.Now().UnixNano(),
		nonce:        0,
		previousHash: b.Hash(),
		transactions: []*Transaction{},
	}
	c := []*Block{appendBlock}
	got := NewBlockchain("test Address", uint16(8000))
	want := &Blockchain{
		transactionPool:   []*Transaction{},
		chain:             c,
		blockchainAddress: "test Address",
	}
	if !reflect.DeepEqual(got.transactionPool, want.transactionPool) {
		t.Errorf("テストに失敗しました。\n got: %#v\nwant: %#v", &got.transactionPool, &want.transactionPool)
	}
	opts := []cmp.Option{
		cmp.AllowUnexported(Block{}),
		cmpopts.IgnoreFields(Block{}, "previousHash", "transactions"),
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
	transactions := []*Transaction(nil)
	got := bc.CreateBlock(nonce, previousHash)
	want := &Block{
		timestamp:    clock.Now().UnixNano(),
		nonce:        nonce,
		previousHash: previousHash,
		transactions: transactions,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("テストに失敗しました。\n got: %#v\nwant: %#v", got, want)
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
		transactions: []*Transaction{},
	}
	b2 := &Block{
		timestamp:    clock.Now().UnixNano() + 1,
		nonce:        2,
		previousHash: b1.Hash(),
		transactions: []*Transaction{},
	}
	b3 := &Block{
		timestamp:    clock.Now().UnixNano() + 2,
		nonce:        3,
		previousHash: b2.Hash(),
		transactions: []*Transaction{},
	}

	for _, test := range []struct {
		title       string
		input_tp    []*Transaction
		input_chain []*Block
		output      *Block
	}{
		{"b1変数のブロックが戻る", []*Transaction{}, []*Block{b1}, b1},
		{"b2変数のブロックが戻る", []*Transaction{}, []*Block{b1, b2}, b2},
		{"b3変数のブロックが戻る", []*Transaction{}, []*Block{b1, b2, b3}, b3},
		{"b2変数のブロックが戻る", []*Transaction{}, []*Block{b3, b1, b2}, b2},
	} {
		t.Run(test.title, func(t *testing.T) {
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
		transactions: []*Transaction{},
	}
	b2 := &Block{
		timestamp:    clock.Now().UnixNano() + 1,
		nonce:        2,
		previousHash: b1.Hash(),
		transactions: []*Transaction{},
	}
	var want string

	for _, test := range []struct {
		title       string
		input_chain []*Block
		output      string
	}{
		{"ターミナルにチェーンの内容が出力される", []*Block{b1, b2}, want},
	} {
		t.Run(test.title, func(t *testing.T) {
			bc := &Blockchain{
				transactionPool: []*Transaction{},
				chain:           test.input_chain,
			}
			for i, b := range bc.chain {
				want += fmt.Sprintf("%s Chain %d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
				want += fmt.Sprintf("nonce        :%d\n", b.nonce)
				want += fmt.Sprintf("previousHash :%x\n", b.previousHash)
				want += fmt.Sprintf("timestamp    :%d\n", b.timestamp)
				if i == len(bc.chain)-1 {
					want += fmt.Sprintf("%s", strings.Repeat("*", 50))
				} else {
					want += fmt.Sprintf("%s\n", strings.Repeat("*", 50))
				}
			}
			got := utils.ExtractStdout(t, bc.Print)
			opts := []cmp.Option{
				cmp.AllowUnexported(Block{}),
				cmpopts.IgnoreFields(Block{}, "previousHash", "transactions"),
			}
			if diff := cmp.Diff(got, want, opts...); diff != "" {
				t.Errorf("テストに失敗しました。\n (-got +want):\n%s", diff)
			}
		})
	}
}
func TestNewTransaction(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		title           string
		input_sender    string
		input_recipient string
		input_value     float32
		output          *Transaction
	}{
		{"引数に応じたトランザクションが作成される", "senderA", "recipientB", 100.05, &Transaction{"senderA", "recipientB", 100.05}},
		{"引数に応じたトランザクションが作成される", "senderA", "recipientB", 0.00064, &Transaction{"senderA", "recipientB", 0.00064}},
		{"引数に応じたトランザクションが作成される", "senderA", "recipientB", -0.05, &Transaction{"senderA", "recipientB", -0.05}},
	} {
		t.Run(test.title, func(t *testing.T) {
			got := NewTransaction(test.input_sender, test.input_recipient, test.input_value)
			want := test.output
			opts := []cmp.Option{
				cmp.AllowUnexported(Transaction{}),
			}
			if diff := cmp.Diff(got, want, opts...); diff != "" {
				t.Errorf("テストに失敗しました。\n (-got +want):\n%s", diff)
			}
		})
	}
}

func TestAddTransaction(t *testing.T) {
	t.Parallel()

	walletA := NewWallet()
	walletB := NewWallet()
	wtr := NewWalletTransaction(walletA.PrivateKey(), walletA.PublicKey(), walletA.BlockchainAddress(), walletB.BlockchainAddress(), 1.0)

	pr, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	r := big.NewInt(int64(mathrand.Intn(99999)))
	s := big.NewInt(int64(mathrand.Intn(99999)))

	for _, test := range []struct {
		title                   string
		input_sender_public_key *ecdsa.PublicKey
		input_signature         *Signature
		input_sender            string
		input_recipient         string
		input_value             float32
		output                  bool
	}{
		{"送信者がMINING_SENDERの場合", nil, nil, MINING_SENDER, walletB.BlockchainAddress(), 3.05646, true},
		{"送信者がMINING_SENDER以外の場合", walletA.PublicKey(), wtr.GenerateSignature(), walletA.BlockchainAddress(), walletB.BlockchainAddress(), 1.0, true},
		{"transactionとvalueが異なる場合", walletA.PublicKey(), wtr.GenerateSignature(), walletA.BlockchainAddress(), walletB.BlockchainAddress(), 100000, false},
		{"不正なパブリックキーの場合", &pr.PublicKey, wtr.GenerateSignature(), walletA.BlockchainAddress(), walletB.BlockchainAddress(), 1.0, false},
		{"不正な署名の場合", walletA.PublicKey(), &Signature{r, s}, walletA.BlockchainAddress(), walletB.BlockchainAddress(), 1.0, false},
		{"受取側のアドレスが異なる場合", walletA.PublicKey(), wtr.GenerateSignature(), walletA.BlockchainAddress(), "wrong address", 1.0, false},
		{"送信者のアドレスが異なる場合", walletA.PublicKey(), wtr.GenerateSignature(), "wrong address", walletB.BlockchainAddress(), 1.0, false},
	} {
		t.Run(test.title, func(t *testing.T) {
			w := NewWallet()
			bc := NewBlockchain(w.BlockchainAddress(), uint16(8000))
			got := bc.AddTransaction(
				test.input_sender_public_key,
				test.input_signature,
				test.input_sender,
				test.input_recipient,
				test.input_value,
			)
			want := test.output
			if got != want {
				t.Errorf("テストに失敗しました。\n got: %#v\nwant: %#v", got, want)
			}
		})
	}
}

func TestVerifyTransactionSignature(t *testing.T) {
	t.Parallel()

	walletA := NewWallet()
	walletB := NewWallet()
	wtr := NewWalletTransaction(walletA.PrivateKey(), walletA.PublicKey(), walletA.BlockchainAddress(), walletB.BlockchainAddress(), 1.0)
	tr := NewTransaction(walletA.BlockchainAddress(), walletB.BlockchainAddress(), 1.0)

	for _, test := range []struct {
		title                   string
		input_sender_public_key *ecdsa.PublicKey
		input_signature         *Signature
		input_transaction       *Transaction
		output                  bool
	}{
		{"Verify", walletA.PublicKey(), wtr.GenerateSignature(), tr, true},
		{"Verify", walletB.PublicKey(), wtr.GenerateSignature(), tr, false},
	} {
		t.Run(test.title, func(t *testing.T) {
			w := NewWallet()
			bc := NewBlockchain(w.BlockchainAddress(), uint16(8000))
			got := bc.VerifyTransactionSignature(
				test.input_sender_public_key,
				test.input_signature,
				test.input_transaction,
			)
			want := test.output
			if got != want {
				t.Errorf("テストに失敗しました。\n got: %#v\nwant: %#v", got, want)
			}
		})
	}
}

func TestCopyTransactionPool(t *testing.T) {
	t.Parallel()

	t1 := NewTransaction("senderA", "recipientB", 1.0)
	t2 := NewTransaction("senderB", "recipientC", 2.0)
	t3 := NewTransaction("senderC", "recipientD", 3.0)
	t4 := NewTransaction("senderD", "recipientE", 4.0)
	bc := new(Blockchain)
	bc.transactionPool = []*Transaction{t1, t2, t3, t4}

	got := bc.CopyTransactionPool()
	if !reflect.DeepEqual(got, bc.transactionPool) {
		t.Errorf("テストに失敗しました。 got:%#v want:%#v", got, bc.transactionPool)
	}
}

func TestValidProof(t *testing.T) {
	t.Parallel()

	bc := new(Blockchain)
	got := bc.ValidProof(0, [32]byte{}, []*Transaction{}, MINING_DIFFICULTY)
	if got != false {
		t.Errorf("テストに失敗しました。\n got: %#v\nwant: %#v", got, false)
	}
}

func TestProofOfWork(t *testing.T) {
	t.Parallel()

	bc := new(Blockchain)
	t1 := NewTransaction("", "", 0.0)
	bc.transactionPool = append(bc.transactionPool, t1)
	b := new(Block)
	bc.chain = append(bc.chain, b)
	got := bc.ProofOfWork()
	want := 1592
	if got != want {
		t.Errorf("テストに失敗しました。 got:%#v want:%#v", got, want)
	}
}

func TestMining(t *testing.T) {
	t.Parallel()

	bc := new(Blockchain)
	t1 := NewTransaction("", "", 0.0)
	bc.transactionPool = append(bc.transactionPool, t1)
	b := new(Block)
	bc.chain = append(bc.chain, b)
	got := bc.Mining()
	want := true
	if got != want {
		t.Errorf("テストに失敗しました。 got:%#v want:%#v", got, want)
	}
}

func TestCalculateTotalAmount(t *testing.T) {
	t.Parallel()

	walletM := NewWallet()
	walletA := NewWallet()
	walletB := NewWallet()

	bc := NewBlockchain(walletM.BlockchainAddress(), uint16(8000))

	value1 := random(0.0, 5.0)
	wtr1 := NewWalletTransaction(walletA.PrivateKey(), walletA.PublicKey(), walletA.BlockchainAddress(), walletB.BlockchainAddress(), value1)
	bc.AddTransaction(walletA.PublicKey(), wtr1.GenerateSignature(), walletA.BlockchainAddress(), walletB.BlockchainAddress(), value1)

	value2 := random(0.0, 10.0)
	wtr2 := NewWalletTransaction(walletB.PrivateKey(), walletB.PublicKey(), walletB.BlockchainAddress(), walletA.BlockchainAddress(), value2)
	bc.AddTransaction(walletB.PublicKey(), wtr2.GenerateSignature(), walletB.BlockchainAddress(), walletA.BlockchainAddress(), value2)

	value3 := random(0.0, 15.0)
	wtr3 := NewWalletTransaction(walletA.PrivateKey(), walletA.PublicKey(), walletA.BlockchainAddress(), walletB.BlockchainAddress(), value3)
	bc.AddTransaction(walletA.PublicKey(), wtr3.GenerateSignature(), walletA.BlockchainAddress(), walletB.BlockchainAddress(), value3)

	value4 := random(0.0, 20.0)
	wtr4 := NewWalletTransaction(walletB.PrivateKey(), walletB.PublicKey(), walletB.BlockchainAddress(), walletA.BlockchainAddress(), value4)
	bc.AddTransaction(walletB.PublicKey(), wtr4.GenerateSignature(), walletB.BlockchainAddress(), walletA.BlockchainAddress(), value4)

	bc.Mining()

	for _, test := range []struct {
		title  string
		input  string
		output float32
	}{
		{"walletAの取引の合計", walletA.BlockchainAddress(), (-value1 + value2 + -value3 + value4)},
		{"walletBの取引の合計", walletB.BlockchainAddress(), (value1 - value2 + value3 - value4)},
	} {
		t.Run(test.title, func(t *testing.T) {
			got := bc.CalculateTotalAmount(test.input)
			want := test.output
			if got != want {
				t.Errorf("テストに失敗しました。\n got: %#v\nwant: %#v", got, want)
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

func random(min, max float32) float32 {
	mathrand.Seed(clock.Now().UnixNano() + int64(mathrand.Intn(99999)))
	return mathrand.Float32()*(max-min) + min
}
