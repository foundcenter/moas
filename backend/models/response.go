package models

type SearchResultAndErrors struct {
	SearchResult []SearchResult `json:"search_result"`
	SearchError  SearchError    `json:"search_error"`
}

type SearchResult struct {
	AccountID   string `json:"account_id"`
	Service     string `json:"service"`
	Resource    string `json:"resource"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
	ID          int    `json:"-"`
}

type SearchError struct {
	AccountID   string      `json:"account_id"`
	AccountType string      `json:"account_type"`
	Error       interface{} `json:"error"`
}

type SearchErrors struct {
	SearchError []SearchError `json:"errors"`
}
