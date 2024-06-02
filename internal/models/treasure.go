package models

import (
	"gorm.io/gorm"
)

type Treasure struct {
	gorm.Model
	ID    int
	Name  string
	Value int
}
