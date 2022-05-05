package models

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewWallet(t *testing.T) {
	t.Parallel()
	setup()

	got := NewWallet()
	want := &Wallet{}
	opts := []cmp.Option{
		cmp.AllowUnexported(Wallet{}),
		cmpopts.IgnoreFields(Wallet{}, "privateKey", "publicKey", "blockchainAddress"),
	}
	if diff := cmp.Diff(got, want, opts...); diff != "" {
		t.Errorf("テストに失敗しました。\n (-got +want):\n%s", diff)
	}
}

func TestPrivateKey(t *testing.T) {
	t.Parallel()
	w := NewWallet()
	got := w.PrivateKey()
	want := w.privateKey
	if got != want {
		t.Errorf("テストに失敗しました。 got:%#v want:%#v", got, want)
	}
}
func TestPrivateKeyStr(t *testing.T) {
	t.Parallel()
	w := NewWallet()
	got := w.PrivateKeyStr()
	want := fmt.Sprintf("%x", w.privateKey.D.Bytes())
	if got != want {
		t.Errorf("テストに失敗しました。 got:%#v want:%#v", got, want)
	}
}

func TestPublicKey(t *testing.T) {
	t.Parallel()
	w := NewWallet()
	got := w.PublicKey()
	want := w.publicKey
	if got != want {
		t.Errorf("テストに失敗しました。 got:%#v want:%#v", got, want)
	}
}

func TestPublicKeyStr(t *testing.T) {
	t.Parallel()
	w := NewWallet()
	got := w.PublicKeyStr()
	want := fmt.Sprintf("%x%x", w.publicKey.X.Bytes(), w.publicKey.Y.Bytes())
	if got != want {
		t.Errorf("テストに失敗しました。 got:%#v want:%#v", got, want)
	}
}
