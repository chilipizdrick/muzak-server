package api

type ErrorResponse struct {
	Error Error `json:"error"`
}

type Error struct {
	Status  uint   `json:"status"`
	Message string `json:"message"`
}
