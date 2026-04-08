package dto

type PostRequest struct {
	URL        string `json:"url"`
	CustomCode string `json:"custom_code"`
}
