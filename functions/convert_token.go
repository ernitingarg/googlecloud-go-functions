package functions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"soteria-functions/client/exchange"
	"soteria-functions/entity"
	"soteria-functions/env"
	"soteria-functions/infra/firestore"
	"soteria-functions/logger"
	"soteria-functions/usecases/balance"

	"github.com/pkg/errors"
)

var convertUsecase *balance.ConvertUsecase

func init() {
	client := firestore.NewFirestoreClient(env.EnvVars.GCP.ProjectID)
	okexClient := exchange.NewOkexClient()
	convertUsecase = balance.NewConvertUsecase(client, okexClient)
}

type convertTokenBody struct {
	FromAmount int64            `json:"from_amount"`
	ToAmount   int64            `json:"to_amount"`
	FromType   entity.AssetType `json:"from_type"`
	ToType     entity.AssetType `json:"to_type"`
}

func (b *convertTokenBody) validate() error {
	if b.FromAmount == 0 {
		return errors.New("FromAmount is zero")
	}

	if b.ToAmount == 0 {
		return errors.New("ToAmount is zero")
	}

	if err := b.FromType.Validate(); err != nil {
		return errors.Wrap(err, "FromType is invalid")
	}

	if err := b.ToType.Validate(); err != nil {
		return errors.Wrap(err, "ToType is invalid")
	}

	return nil
}

// ConvertToken btc address
func ConvertToken(w http.ResponseWriter, r *http.Request) {
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

	bodyStr, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeErrorResponse(w, http.StatusInternalServerError, entity.SERVER_ERROR, "failed read request body", err)
		logger.Log.Error(errors.Wrap(err, "failed read request body"))
		return
	}

	var body convertTokenBody
	err = json.Unmarshal(bodyStr, &body)
	logger.Log.Info(bodyStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeErrorResponse(w, http.StatusInternalServerError, entity.SERVER_ERROR, "failed request body decode", err)
		return
	}

	if err := body.validate(); err != nil {
		msg := fmt.Sprintf("invalid body. %+v", body)
		w.WriteHeader(http.StatusBadRequest)
		writeErrorResponse(w, http.StatusBadRequest, entity.VALIDATION_ERROR, msg, err)
		return
	}

	if err = convertUsecase.ValidatePrice(ctx, body.FromAmount, body.ToAmount, body.FromType, body.ToType); err != nil {
		msg := "invalid price."
		w.WriteHeader(http.StatusBadRequest)
		writeErrorResponse(w, http.StatusBadRequest, entity.VALIDATION_ERROR, msg, err)
		return
	}

	err = convertUsecase.Convert(ctx, uid, body.FromAmount, body.ToAmount, body.FromType, body.ToType)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeErrorResponse(w, http.StatusInternalServerError, entity.SERVER_ERROR, "failed convert", err)
		return
	}

	logger.Log.Info(fmt.Sprintf("success convert token. body: %+v", body))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"status\": \"OK\"}"))
}
