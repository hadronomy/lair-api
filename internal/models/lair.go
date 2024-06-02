package models

import (
	"net/http"
)

type Lair struct {
	Model
	Name      string     `json:"name"`
	Owner     string     `json:"owner"`
	Private   bool       `json:"private"`
	Treasures []Treasure `json:"treasures" gorm:"many2many:lair_treasures;"`
	Minions   []Minion   `json:"minions" gorm:"many2many:lair_minions;"`
}

type LairRequest struct {
	Name    string `json:"name" required:"false"`
	Owner   string `json:"owner" required:"false"`
	Private bool   `json:"private" required:"false"`
}

func (l *Lair) Bind(r *http.Request) error {
	return nil
}
