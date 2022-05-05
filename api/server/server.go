package server

import (
	"app/models"
	"io"
	"log"
	"net/http"
	"strconv"
)

var cache map[string]*models.Blockchain = make(map[string]*models.Blockchain)

type BlockchainServer struct {
	port uint16
}

func NewBlockchainServer(port uint16) *BlockchainServer {
	return &BlockchainServer{port}
}

func (bcs *BlockchainServer) Port() uint16 {
	return bcs.port
}

func (bcs *BlockchainServer) GetBlockchain() *models.Blockchain {
	bc, isExists := cache["blockchain"]
	if isExists {
		return bc
	}
	newWallet := models.NewWallet()
	bc = models.NewBlockchain(newWallet.BlockchainAddress(), bcs.Port())
	cache["blockchain"] = bc
	return bc
}

func (bcs *BlockchainServer) GetChain(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		bc := bcs.GetBlockchain()
		m, _ := bc.MarshalJSON()
		io.WriteString(w, string(m))
	default:
		log.Printf("ERROR: Invalid HTTP Method")
	}

}

func (bcs *BlockchainServer) Run() {
	http.HandleFunc("/", bcs.GetChain)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(bcs.port)), nil))
}
