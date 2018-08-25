package main

import "github.com/jinzhu/gorm"

// DB TYPES

type Status struct {
	gorm.Model
	One      string `gorm:"default:'one'"`
	HistFrom int64  `gorm:"default:'0'"`
}
