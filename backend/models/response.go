package models

type SearchResult struct {
	AccountID   string `json:"account_id"`
	Service     string `json:"service"`
	Resource    string `json:"resource"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
}
