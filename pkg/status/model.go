package status

type Status struct {
	Name string `json:"name"`
	IsUp bool   `json:"isUp"`
	IP   string `json:"ip"`
}
