package httperror

type HTTPError struct {
	Message interface{} `json:"message"`
}
