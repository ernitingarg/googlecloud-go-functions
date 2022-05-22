package soteriafunctions

import (
	"context"
	"net/http"
	"soteria-functions/functions"
	"soteria-functions/logger"
)

// CreateBtcAccount verify user & re PUT
func CreateBtcAccount(w http.ResponseWriter, r *http.Request) {
	defer logger.Log.Flush()
	result := corsEnabledFunctionAuth(w, r)
	if !result {
		logger.Log.Info("[CreateBtcAccount] result is false")
		return
	}
	functions.CreateBtcAccount(w, r)
}

// CreateEthAccount verify user & re PUT
func CreateEthAccount(w http.ResponseWriter, r *http.Request) {
	defer logger.Log.Flush()
	result := corsEnabledFunctionAuth(w, r)
	if !result {
		logger.Log.Info("[CreateEthAccount] result is false")
		return
	}
	functions.CreateEthAccount(w, r)
}

// AddBtcBalance verify user & re PUT
func AddBtcBalance(w http.ResponseWriter, r *http.Request) {
	defer logger.Log.Flush()
	// result := corsEnabledFunctionAuth(w, r)
	// if !result {
	// 	logger.Log.Info("[AddBtcBalance] result is false")
	// 	return
	// }
	functions.AddBtcBalance(w, r)
}

// AddUsdsBalance verify user & re PUT
func AddUsdsBalance(w http.ResponseWriter, r *http.Request) {
	defer logger.Log.Flush()
	// result := corsEnabledFunctionAuth(w, r)
	// if !result {
	// 	logger.Log.Info("[AddBtcBalance] result is false")
	// 	return
	// }
	functions.AddUsdsBalance(w, r)
}

// AddUsdsBalance verify user & re PUT
func ConvertToken(w http.ResponseWriter, r *http.Request) {
	defer logger.Log.Flush()
	result := corsEnabledFunctionAuth(w, r)
	if !result {
		logger.Log.Info("[ConvertToken] result is false")
		return
	}
	functions.ConvertToken(w, r)
}

// PubSubMessage is the payload of a Pub/Sub event.
type pubSubMessage struct {
	Data []byte `json:"data"`
}

// UpdatePriceHistory
func UpdatePriceHistory(ctx context.Context, m pubSubMessage) error {
	defer logger.Log.Flush()
	return functions.UpdateBtcPrice(ctx)
}

// corsEnabledFunctionAuth https://cloud.google.com/functions/docs/writing/http?hl=ja#authentication_and_cors
// preflightの時はheaderをセットするだけで本実行してはいけない
func corsEnabledFunctionAuth(w http.ResponseWriter, r *http.Request) bool {
	// Set CORS headers for the preflight request
	// Originヘッダーはリクエスト元のホスト名を表している
	host := r.Header.Get("Origin")
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "POST, PUT")
		w.Header().Set("Access-Control-Allow-Origin", host)
		w.Header().Set("Access-Control-Max-Age", "3600")
		return false
	}
	// Set CORS headers for the main request.
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", host)
	return true
}
