package entity

import "github.com/pkg/errors"

type PriceHistory struct {
	CurrencyPair string  `json:"currency_pair"`
	Rate         float64 `json:"rate"`
	Source       string  `json:"source"`
	Timestamp    uint32  `json:"timestamp"`
}

func NewPriceHistory(currencyPair string, rate float64, source string, timestamp uint32) *PriceHistory {
	return &PriceHistory{CurrencyPair: currencyPair, Rate: rate, Source: source, Timestamp: timestamp}
}

func (p *PriceHistory) Validate() error {
	if p.CurrencyPair == "" {
		return errors.New("currency pair is empty")
	}
	if p.Rate <= 0 {
		return errors.New("rate is less than zero")
	}
	if p.Source == "" {
		return errors.New("source is empty")
	}
	if p.Timestamp <= 0 {
		return errors.New("timestamp is less than zero")
	}
	return nil
}
