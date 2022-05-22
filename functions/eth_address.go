package functions

import (
	"errors"
	"fmt"
	"net/http"
	"soteria-functions/entity"
	"soteria-functions/env"
	"soteria-functions/infra/firestore"
	"soteria-functions/logger"
	"soteria-functions/usecases/auth"
	"soteria-functions/usecases/eth"
)

var ethUsecase *eth.EthUsecase
var authUsecase *auth.AuthUsecase

func init() {
	client := firestore.NewFirestoreClient(env.EnvVars.GCP.ProjectID)
	ethUsecase = eth.NewEthUsecase(client)
	authUsecase = auth.NewAuthUsecase()
}

// CreateEthAccount eth address
func CreateEthAccount(w http.ResponseWriter, r *http.Request) {
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

	// create eth account
	account, err := ethUsecase.CreateEthAccount(uid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeErrorResponse(w, http.StatusInternalServerError, entity.SERVER_ERROR, "failed CreateEthAccount", err)
		return
	}

	err = ethUsecase.SaveEthAccount(ctx, account)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeErrorResponse(w, http.StatusInternalServerError, entity.SERVER_ERROR, "failed SaveEthAccount", err)
		return
	}

	logger.Log.Info(fmt.Sprintf("success create eth address. uid: %+v address: %+v", account.UserID, account.Address))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"status\": \"OK\"}"))
}
