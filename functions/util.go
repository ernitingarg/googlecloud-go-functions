package functions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"soteria-functions/entity"
	"soteria-functions/logger"
)

func writeErrorResponse(w http.ResponseWriter, statusCode int, errcode entity.ErrorCode, msg string, targetError error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	errRes, err := json.Marshal(entity.NewErrorResponse(errcode, msg))
	if err != nil {
		errRes = []byte(fmt.Sprintf("{\"id\":%d,error:\"unknown\"}", entity.SERVER_ERROR))
		logger.Log.Error(fmt.Sprintf("error occured: %+v", err))
	}
	logger.Log.Error(fmt.Sprintf("return error response. msg: %s. caused: %+v", msg, targetError))
	w.Write(errRes)
	return
}
