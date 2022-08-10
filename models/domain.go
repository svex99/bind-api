package models

type Domain struct {
	Id         uint64      `gorm:"primaryKey"`
	Name       string      `json:"name" binding:"required,min=1" gorm:"unique"`
	NameServer string      `json:"name_server" binding:"required,min=1"`
	NSIp       string      `json:"nsip" binding:"required,ip"`
	Ttl        string      `json:"ttl" binding:"min=1"`
	Subdomains []Subdomain `gorm:"constraint:notOnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type DomainInfo struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}

type UpdateDomainForm struct {
	Name       string `json:"name" binding:"omitempty,min=1"`
	NameServer string `json:"name_server" binding:"omitempty,min=1"`
	NSIp       string `json:"nsip" binding:"omitempty,ip"`
	Ttl        string `json:"ttl" binding:"omitempty,min=1"`
}
