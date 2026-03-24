package model

// JSONRequest - формат JSON-запроса к серверу.
type JSONRequest struct {
	URL string `json:"url"`
}

// JSONResponse - формат JSON-ответа от сервера.
type JSONResponse struct {
	Result string `json:"result"`
}
