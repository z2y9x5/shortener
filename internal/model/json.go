package model

// JSONRequest - формат JSON-запроса к серверу.
type JSONRequest struct {
	URL string `json:"url"`
}

// JSONResponse - формат JSON-ответа от сервера.
type JSONResponse struct {
	Result string `json:"result"`
}

// FileRecord описывает поля записи в файле.
type FileRecord struct {
	Uuid        string `json:"uuid"`
	ShortUrl    string `json:"short_url"`
	OriginalUrl string `json:"original_url"`
}
