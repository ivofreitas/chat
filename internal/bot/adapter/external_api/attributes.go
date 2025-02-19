package client

import (
	"fmt"
)

type Stock struct {
	Symbol string
	Date   string
	Time   string
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
}

func (s *Stock) String() string {
	return fmt.Sprintf("%s quote is $%.2f per share", s.Symbol, s.Close)
}
