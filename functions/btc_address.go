package functions

import (
	"fmt"
	"net/http"
	"soteria-functions/entity"
	"soteria-functions/env"
	"soteria-functions/infra/firestore"
	"soteria-functions/logger"
	"soteria-functions/usecases/auth"
	"soteria-functions/usecases/btc"

	"github.com/pkg/errors"
)

var btcUsecase *btc.BtcUsecase

func init() {
	client := firestore.NewFirestoreClient(env.EnvVars.GCP.ProjectID)
	btcUsecase = btc.NewBtcUsecase(client)
	authUsecase = auth.NewAuthUsecase()
}

// CreateBtcAccount btc address
func CreateBtcAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != "POST" {
		msg := fmt.Sprintf("no post method: %+v", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		writeErrorResponse(w, http.StatusMethodNotAllowed, entity.VALIDATION_ERROR, msg, errors.New(msg))
		return
	}

	uid, err := authUsecase.GetUID(ctx, r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		writeErrorResponse(w, http.StatusUnauthorized, entity.AUTH_ERROR, "failed get uid", err)
		return
	}

	// create btc account
	account, err := btcUsecase.CreateBtcAccount(uid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeErrorResponse(w, http.StatusInternalServerError, entity.SERVER_ERROR, "failed CreateBtcAccount", err)
		return
	}

	err = btcUsecase.SaveBtcAccount(ctx, account)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeErrorResponse(w, http.StatusInternalServerError, entity.SERVER_ERROR, "failed SaveBtcAccount", err)
		return
	}

	logger.Log.Info(fmt.Sprintf("success create btc address. uid: %+v address: %+v", account.UserID, account.Address))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"status\": \"OK\"}"))
}
