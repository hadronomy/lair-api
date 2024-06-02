package models

import (
	"gorm.io/gorm"
)

type Minion struct {
	gorm.Model
	ID    int
	Name  string
	Level int
}
