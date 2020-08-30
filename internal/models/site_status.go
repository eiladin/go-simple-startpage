package models

type SiteStatus struct {
	ID   uint   `json:"id"`
	IsUp bool   `json:"isUp"`
	IP   string `json:"ip"`
}
