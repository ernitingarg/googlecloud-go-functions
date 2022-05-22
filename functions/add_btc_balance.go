package functions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"soteria-functions/entity"
	"soteria-functions/env"
	"soteria-functions/infra/firestore"
	"soteria-functions/logger"
	"soteria-functions/usecases/balance"

	"github.com/pkg/errors"
)

var btcBalanceUsecase *balance.BtcBalanceUsecase

func init() {
	client := firestore.NewFirestoreClient(env.EnvVars.GCP.ProjectID)
	btcBalanceUsecase = balance.NewBtcBalanceUsecase(client)
}

type addBtcBalanceBody struct {
	Address string `json:"address"`
	Amount  int64  `json:"amount"`
	Txid    string `json:"txid"`
}

// AddBtcBalance btc address
func AddBtcBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != "POST" {
		msg := fmt.Sprintf("no post method: %+v", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		writeErrorResponse(w, http.StatusMethodNotAllowed, entity.VALIDATION_ERROR, msg, errors.New(msg))
		return
	}

	bodyStr, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeErrorResponse(w, http.StatusInternalServerError, entity.SERVER_ERROR, "failed read request body", err)
		return
	}

	var body addBtcBalanceBody
	err = json.Unmarshal(bodyStr, &body)
	logger.Log.Info(bodyStr)
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

	err = btcBalanceUsecase.AddBalance(ctx, body.Address, body.Txid, body.Amount)
	if err != nil {
		msg := fmt.Sprintf("failed btcBalanceUsecase.AddBalance. address: address:%s, txid: %s, amount: %d", body.Address, body.Txid, body.Amount)
		w.WriteHeader(http.StatusInternalServerError)
		writeErrorResponse(w, http.StatusInternalServerError, entity.SERVER_ERROR, msg, err)
		return
	}

	logger.Log.Info(fmt.Sprintf("success add btc. address: address:%s, txid: %s, amount: %d", body.Address, body.Txid, body.Amount))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"status\": \"OK\"}"))
}
