package models

import "time"

type Statistic struct {
	Avg float64 `json:"avg"`
	Max float64 `json:"max"`
	Min float64 `json:"min"`
}

type Timeseries struct {
	Timestamp   time.Time `json:"ts"`
	Cpu         float32   `json:"cpu"`
	Concurrency uint32    `json:"concurrency"`
}

type Message struct {
	Error string `json:"error"`
}
