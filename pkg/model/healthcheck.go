package model

type Healthcheck struct {
	Status string            `json:"status"`
	Errors HealthcheckErrors `json:"errors"`
}

type HealthcheckErrors struct {
	Database string `json:"database"`
}
