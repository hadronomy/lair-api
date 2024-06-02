package models

type Treasure struct {
	Model
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Value int    `json:"value"`
}
