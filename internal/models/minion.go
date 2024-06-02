package models

type Minion struct {
	Model
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Level int    `json:"level"`
}
