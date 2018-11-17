package main

import "time"

// NATS TYPES

// COMMON

type response struct {
	Err error
}

// HISTORY

type historyQuery struct {
	APIKey string
	Query  string
}

type historyResults struct {
	HistoryEntries []HistoryEntry
}

type historyNew struct {
	APIKey     string
	Entry      string
	UserAtHost string
	CreatedAt  time.Time
}

type historyNewResp struct {
	Error error
}

// PROMPT

type promptQuery struct {
	APIKey string
	Query  string
}

type promptResults struct {
	PromptEntries []Prompt
}

// ACCOUNT
// create

type accountCreateReq struct {
	Email    string
	Password string
}

type accountCreateRes struct {
	Status  bool
	Message string
}

// login

type accountLoginReq struct {
	Email    string
	Password string
}

type accountLoginRes struct {
	Status  bool
	Message string
	APIKey  string
}
