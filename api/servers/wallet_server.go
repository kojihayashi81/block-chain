package servers

import (
	"app/middleware"
	"app/models"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
)

type WalletServer struct {
	port    uint16
	gateway string
}

func NewWalletServer(port uint16, gateway string) *WalletServer {
	return &WalletServer{port, gateway}
}

func (ws *WalletServer) Port() uint16 {
	return ws.port
}

func (ws *WalletServer) Gateway() string {
	return ws.gateway
}

func (ws *WalletServer) GetWallet(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		middleware.SetupResponseWriter(&w)
		wt := models.NewWallet()
		m, _ := wt.MarshalJSON()
		io.WriteString(w, string(m))
	default:
		log.Printf("ERROR: Invalid HTTP Method")
	}
}

func (ws *WalletServer) CreateTransaction(w http.ResponseWriter, req *http.Request) {
	middleware.SetupResponseWriter(&w)
	switch req.Method {
	case http.MethodPost:
		decoder := json.NewDecoder(req.Body)
		var t models.TransactionRequest
		err := decoder.Decode(&t)

		if err != nil {
			log.Printf("ERROR: %v", err)
			io.WriteString(w, "fail")
			return
		}

		if !t.ValidateTransactionRequest() {
			log.Printf("ERROR: missing field(s)")
			io.WriteString(w, "missing field(s)")
			return
		}

		publicKey := models.PublicKeyFromString(*t.SenderPublicKey)
		privateKey := models.PrivateKeyFromString(*t.SenderPrivateKey, publicKey)
		value, err := strconv.ParseFloat(*t.Value, 32)
		if err != nil {
			log.Panicln("ERROR: parse error")
			io.WriteString(w, "fail")
		}
		value32 := float32(value)
		w.Header().Add("Content-Type", "application/json")

		transaction := models.NewWalletTransaction(
			privateKey,
			publicKey,
			*t.SenderBlockchainAddress,
			*t.RecipientBlockchainAddress,
			value32,
		)
		signature := transaction.GenerateSignature()
		signatureStr := signature.String()

		btr := &models.BlockchainTransactionRequest{
			SenderPublicKey:            t.SenderPublicKey,
			SenderBlockchainAddress:    t.SenderBlockchainAddress,
			RecipientBlockchainAddress: t.RecipientBlockchainAddress,
			Value:                      &value32,
			Signature:                  &signatureStr,
		}
		m, _ := json.Marshal(btr)
		buf := bytes.NewBuffer(m)

		resp, err := http.Post("http://0.0.0.0:8000/transactions", "application/json", buf)
		if err != nil {
			log.Printf("ERROR: %v", err)
			io.WriteString(w, "fail")
			return
		}
		if resp.StatusCode == 201 {
			io.WriteString(w, "success")
			return
		}
		io.WriteString(w, "fail")
	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)
		break
	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("ERROR: Invalid HTTP Method")
	}
}

func (ws *WalletServer) WalletAmount(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		middleware.SetupResponseWriter(&w)
		blockchainAddress := req.URL.Query().Get("blockchain_address")
		endpoint := "http://0.0.0.0:8000/amount"
		client := &http.Client{}
		bcsReq, _ := http.NewRequest("GET", endpoint, nil)
		q := bcsReq.URL.Query()
		q.Add("blockchain_address", blockchainAddress)
		bcsReq.URL.RawQuery = q.Encode()

		bcsResp, err := client.Do(bcsReq)
		if err != nil {
			log.Printf("ERROR: %v", err)
			io.WriteString(w, "fail")
		}

		if bcsResp.StatusCode == 200 {
			w.Header().Add("Content-Type", "application/json")
			decoder := json.NewDecoder(bcsResp.Body)
			var bar models.AmountRequest
			err := decoder.Decode(&bar)
			if err != nil {
				log.Printf("ERROR: %v", err)
				io.WriteString(w, "fail")
				return
			}
			m, _ := json.Marshal(struct {
				Message string  `json:"message"`
				Amount  float32 `json:"amount"`
			}{
				Message: "success",
				Amount:  bar.Amount,
			})
			io.WriteString(w, string(m[:]))
		} else {
			io.WriteString(w, "fail")
		}
	default:
		log.Printf("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}
}
