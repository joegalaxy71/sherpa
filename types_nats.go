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
	HistoryEntries []historyEntry
}

type historyEntry struct {
	Account    int64
	Entry      string
	UserAtHost string
	CreatedAt  time.Time
}

type historyNewResp struct {
	Error error
}
