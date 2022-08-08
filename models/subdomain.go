package models

type Subdomain struct {
	Name string `json:"name" binding:"required"`
	Ip   string `json:"ip" binding:"required"`
}
