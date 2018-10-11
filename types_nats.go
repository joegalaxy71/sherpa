package main

import "time"

// NATS TYPES

// COMMON

type response struct {
	Err error
}

// HISTORY

type historyQuery struct {
	Query string
}

type historyResults struct {
	HistoryEntries []HistoryEntry
}

type historyNew struct {
	Account    string
	Entry      string
	UserAtHost string
	CreatedAt  time.Time
}

type historyNewResp struct {
	Error error
}

// ACCOUNT
// create

type accountCreateReq struct {
	Email    string
	Password string
}
