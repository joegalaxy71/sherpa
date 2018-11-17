package main

import "github.com/jinzhu/gorm"

// DB TYPES

type Status struct {
	gorm.Model
	One      string `gorm:"default:'one'"`
	HistFrom int64  `gorm:"default:'0'"`
}

type HistoryEntry struct {
	gorm.Model
	Account    string
	Entry      string
	UserAtHost string
}

type Prompt struct {
	gorm.Model
	Name     string
	Tag      string
	Preview  string
	Sequence string
}
