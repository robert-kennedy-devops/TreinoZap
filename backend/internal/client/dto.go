package client

type CreateRequest struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Goal  string `json:"goal"`
	Notes string `json:"notes"`
}

type UpdateRequest struct {
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	Status string `json:"status"`
	Goal   string `json:"goal"`
	Notes  string `json:"notes"`
}

type ListResponse struct {
	Data       []Client `json:"data"`
	Total      int      `json:"total"`
	Page       int      `json:"page"`
	PageSize   int      `json:"page_size"`
	TotalPages int      `json:"total_pages"`
}
