package exchange

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	ftxBaseURL string = "https://ftx.com/api"
)

type FtxClient struct {
	client *http.Client
}

func NewFtxClient() *FtxClient {
	client := new(http.Client)
	return &FtxClient{client: client}
}

type result struct {
	Ask                   float64     `json:"ask"`
	BaseCurrency          interface{} `json:"baseCurrency"`
	Bid                   float64     `json:"bid"`
	Change1H              float64     `json:"change1h"`
	Change24H             float64     `json:"change24h"`
	ChangeBod             float64     `json:"changeBod"`
	Enabled               bool        `json:"enabled"`
	HighLeverageFeeExempt bool        `json:"highLeverageFeeExempt"`
	Last                  float64     `json:"last"`
	MinProvideSize        float64     `json:"minProvideSize"`
	Name                  string      `json:"name"`
	PostOnly              bool        `json:"postOnly"`
	Price                 float64     `json:"price"`
	PriceIncrement        float64     `json:"priceIncrement"`
	QuoteCurrency         interface{} `json:"quoteCurrency"`
	QuoteVolume24H        float64     `json:"quoteVolume24h"`
	Restricted            bool        `json:"restricted"`
	SizeIncrement         float64     `json:"sizeIncrement"`
	Type                  string      `json:"type"`
	Underlying            string      `json:"underlying"`
	VolumeUsd24H          float64     `json:"volumeUsd24h"`
}

type response struct {
	Success bool   `json:"success"`
	Result  result `json:"result"`
}

func (o *FtxClient) GetPrice(ctx context.Context, marketName string) (float64, error) {
	endpoint := fmt.Sprintf("%s/markets/%s", ftxBaseURL, marketName)
	req, _ := http.NewRequest("GET", endpoint, nil)
	req.Header.Set("Content-Type", "application/json")
	res, err := o.client.Do(req)
	if err != nil {
		return 0, err
	}

	defer res.Body.Close()
	bytes, err := ioutil.ReadAll(res.Body)
	var data response
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return 0, err
	}

	fmt.Printf("endpoint: %s, data: %+v", endpoint, data)
	if data.Result.Price <= 0 {
		return 0, errors.New("FtxClient:GetPrice: price is less than zero")
	}
	return data.Result.Price, nil
}
