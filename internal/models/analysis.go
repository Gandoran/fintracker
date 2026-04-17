package models

import "time"

type Analysis struct {
	Ticker         []string  `json:"tickers"`
	Sentiment      string    `json:"sentiment"`
	Summary        string    `json:"summary"`
	Impact         string    `json:"impact"`
	ReferenceLinks []string  `json:"reference_links"`
	AnalysisAt     time.Time `json:"analysis_at"`
	Original       Article   `json:"original"`
}
