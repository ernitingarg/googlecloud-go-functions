package exchange

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type InstrumentID string

func (l InstrumentID) String() string {
	return string(l)
}

const (
	BTC_USD     InstrumentID = "BTC-USD"
	okexBaseURL string       = "https://www.okex.com"
)

type OkexClient struct {
	client *http.Client
}

func NewOkexClient() *OkexClient {
	client := new(http.Client)
	return &OkexClient{client: client}
}

type constituents struct {
	Symbol         string `json:"symbol"`
	Original_price string `json:"original_price"`
	Weight         string `json:"weight"`
	Usd_price      string `json:"usd_price"`
	Exchange       string `json:"exchange"`
}

type indexConstituents struct {
	Last         string         `json:"last"`
	Constituents []constituents `json:"constituents"`
}

type indexConstituentsResponse struct {
	Code      int               `json:"code"`
	DetailMsg string            `json:"detailMsg"`
	Msg       string            `json:"msg"`
	Data      indexConstituents `json:"data"`
}

func (o *OkexClient) GetPrice(ctx context.Context, instrumentID InstrumentID) (float64, error) {
	endpoint := fmt.Sprintf("%s/api/index/v3/%s/constituents", okexBaseURL, instrumentID)
	req, _ := http.NewRequest("GET", endpoint, nil)
	req.Header.Set("Content-Type", "application/json")
	res, err := o.client.Do(req)
	if err != nil {
		return 0, err
	}

	defer res.Body.Close()
	bytes, err := ioutil.ReadAll(res.Body)
	var data indexConstituentsResponse
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return 0, err
	}

	fmt.Printf("endpoint: %s, data: %+v", endpoint, data)
	last, err := strconv.ParseFloat(data.Data.Last, 64)
	if err != nil {
		return 0, err
	}
	return last, nil
}
