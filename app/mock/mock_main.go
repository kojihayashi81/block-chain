package mock

import "app"

type AppMock struct {
	app.App
	FakeCreate func(nonce int, previousHash string) *app.Block
}

func (a *AppMock) CreateBlock(nonce int, previousHash string) *app.Block {
	return a.FakeCreate(nonce, previousHash)
}
