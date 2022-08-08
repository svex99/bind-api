package models

type Domain struct {
	Name       string `json:"name" binding:"required,min=1" gorm:"primaryKey"`
	NameServer string `json:"name_server" binding:"required,min=1"`
	NSIp       string `json:"nsip" binding:"required,ip"`
	Ttl        string `json:"ttl" binding:"min=1"`
}

type UpdateDomainForm struct {
	NameServer string `json:"name_server" binding:"min=1"`
	NSIp       string `json:"nsip" binding:"ip"`
	Ttl        string `json:"ttl" binding:"min=1"`
}
