package main

import (
	"github.com/jinzhu/gorm"
)

// MICROSERVER

type microServer struct {
	name    string
	init    func() error
	serve   func(request) response
	cleanup func() error
}

// REQUEST / RESPONSE

type request struct {
	Req string
}

type response struct {
	Res     string
	List    []string
	Prompts []prompt
}

//REMEMBER ALWAYS
// the less identifiers, the best

type prompt struct {
	Name  string
	Value string
}

// DB TYPES

type HistoryEntry struct {
	gorm.Model
	Entry string
}

type Status struct {
	HistSize uint32
}
