package model

type User struct {
	Email string `json:"email"` // unique id
	Name  string `json:"name"`
	Age   int    `json:"age"`

