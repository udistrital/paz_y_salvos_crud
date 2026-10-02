package models

// APIResponse is the response envelope exposed by the CRUD endpoints.
type APIResponse struct {
	Success bool        `json:"Success"`
	Status  int         `json:"Status"`
	Message interface{} `json:"Message"`
	Data    interface{} `json:"Data"`
}
