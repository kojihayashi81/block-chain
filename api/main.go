package main

import (
	"app/servers"
	"flag"
	"log"
	"net/http"
)

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	blockchainServer := http.NewServeMux()
	bcs := servers.NewBlockchainServer(uint16(8000))
	blockchainServer.HandleFunc("/blockchain", bcs.GetChain)
	blockchainServer.HandleFunc("/transactions", bcs.Transactions)
	blockchainServer.HandleFunc("/mine", bcs.Mine)
	blockchainServer.HandleFunc("/mine/start", bcs.StartMine)
	blockchainServer.HandleFunc("/amount", bcs.Amount)

	walletServer := http.NewServeMux()
	gateway := flag.String("gateway", "http://127.0.0.1:5000", "Blockchain Gateway")
	ws := servers.NewWalletServer(uint16(5000), *gateway)
	walletServer.HandleFunc("/wallet", ws.GetWallet)
	walletServer.HandleFunc("/wallet/amount", ws.WalletAmount)
	walletServer.HandleFunc("/transaction", ws.CreateTransaction)

	go func() {
		http.ListenAndServe("0.0.0.0:8000", blockchainServer)
	}()
	http.ListenAndServe("127.0.0.1:5000", walletServer)
}
