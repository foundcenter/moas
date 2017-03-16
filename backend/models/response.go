package models

type ResultResponse struct {
	AccountID   string `json:"account_id"`
	Service     string `json:"service"`
	Resource    string `json:"resource"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
}
