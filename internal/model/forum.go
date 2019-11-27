package model

type Forum struct {
	ID   int64  `json:"-" valid:"int, optional"`
	Name string `json:"name" valid:"utfletter"`
}
