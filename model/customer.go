package model

type Customer struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	PhoneNo string `json:"phone"`
	Address string `json:"address"`
}
