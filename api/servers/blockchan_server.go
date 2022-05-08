package servers

import (
	"app/middleware"
	"app/models"
	"encoding/json"
	"io"
	"log"
	"net/http"
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

func (bcs *BlockchainServer) Transactions(w http.ResponseWriter, req *http.Request) {
	middleware.SetupResponseWriter(&w)
	switch req.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		bc := bcs.GetBlockchain()
		transactions := bc.TransactionPool()
		m, _ := json.Marshal(struct {
			Transactions []*models.Transaction `json:"transactions"`
			Length       int                   `json:"length"`
		}{
			Transactions: transactions,
			Length:       len(transactions),
		})
		io.WriteString(w, string(m[:]))
	case http.MethodPost:
		decoder := json.NewDecoder(req.Body)
		var t models.BlockchainTransactionRequest
		err := decoder.Decode(&t)

		if err != nil {
			log.Printf("ERROR: %v", err)
			io.WriteString(w, "fail")
			return
		}

		if !t.ValidateBlockchainTransactionRequest() {
			log.Printf("ERROR: missing field(s)")
			io.WriteString(w, "missing field(s)")
			return
		}

		publicKey := models.PublicKeyFromString(*t.SenderPublicKey)
		signature := models.SignatureFromString(*t.Signature)
		bc := bcs.GetBlockchain()
		isCreated := bc.CreateTransaction(
			publicKey,
			signature,
			*t.SenderBlockchainAddress,
			*t.RecipientBlockchainAddress,
			*t.Value,
		)

		w.Header().Add("Content-Type", "application/json")
		var m string
		if isCreated {
			w.WriteHeader(http.StatusCreated)
			m = "success"
		} else {
			w.WriteHeader(http.StatusBadRequest)
			m = "fail"
		}
		io.WriteString(w, m)
	default:
		log.Printf("ERROR: Invalid HTTP Method")
	}
}
