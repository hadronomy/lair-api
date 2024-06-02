package models

import (
	"net/http"

	"gorm.io/gorm"
)

type Lair struct {
	gorm.Model
	Name      string
	Owner     string
	Private   bool
	Treasures []Treasure `gorm:"many2many:lair_treasures;"`
	Minions   []Minion   `gorm:"many2many:lair_minions;"`
}

type LairRequest struct {
	Name    string
	Owner   string
	Private bool
}

func (l *LairRequest) Bind(r *http.Request) error {
	return nil
}

func (l *LairRequest) ToLair(id int) Lair {
	return Lair{
		Name:      l.Name,
		Owner:     l.Owner,
		Private:   l.Private,
		Treasures: []Treasure{},
		Minions:   []Minion{},
	}
}

func (l *Lair) Bind(r *http.Request) error {
	return nil
}
