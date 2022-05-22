package functions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"soteria-functions/entity"
	"soteria-functions/env"
	"soteria-functions/infra/firestore"
	"soteria-functions/logger"
	"soteria-functions/usecases/balance"

	"github.com/pkg/errors"
)

var ethBalanceUsecase *balance.EthBalanceUsecase

func init() {
	client := firestore.NewFirestoreClient(env.EnvVars.GCP.ProjectID)
	ethBalanceUsecase = balance.NewEthBalanceUsecase(client)
}

type addUsdsBalanceBody struct {
	Address string `json:"address"`
	Amount  int64  `json:"amount"`
	Txid    string `json:"txid"`
}

// AddUsdsBalance usds address
func AddUsdsBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != "POST" {
		msg := fmt.Sprintf("no post method: %+v", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		writeErrorResponse(w, http.StatusMethodNotAllowed, entity.VALIDATION_ERROR, msg, errors.New(msg))
		return
	}

	var body addUsdsBalanceBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeErrorResponse(w, http.StatusInternalServerError, entity.SERVER_ERROR, "failed request body decode", err)
		return
	}

	if body.Address == "" || body.Amount == 0 || body.Txid == "" {
		msg := fmt.Sprintf("empty body. address: %+v, amount: %+v, txid: %+v", body.Address, body.Amount, body.Txid)
		w.WriteHeader(http.StatusBadRequest)
		writeErrorResponse(w, http.StatusBadRequest, entity.VALIDATION_ERROR, msg, errors.New(msg))
		return
	}

	err = ethBalanceUsecase.AddUsdsBalance(ctx, body.Address, body.Txid, body.Amount)
	if err != nil {
		msg := fmt.Sprintf("failed ethBalanceUsecase.AddUsdsBalance. address: address:%s, txid: %s, amount: %d", body.Address, body.Txid, body.Amount)
		w.WriteHeader(http.StatusInternalServerError)
		writeErrorResponse(w, http.StatusInternalServerError, entity.SERVER_ERROR, msg, err)
		return
	}

	logger.Log.Info(fmt.Sprintf("success add usds. address: address:%s, txid: %s, amount: %d", body.Address, body.Txid, body.Amount))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"status\": \"OK\"}"))
}
