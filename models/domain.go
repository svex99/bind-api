package models

type Domain struct {
	Id         uint64      `json:"id" gorm:"primaryKey"`
	Name       string      `json:"name" binding:"required,min=1" gorm:"unique"`
	NameServer string      `json:"nameServer" binding:"required,min=1"`
	NSIp       string      `json:"nsIp" binding:"required,ip"`
	Ttl        string      `json:"ttl" binding:"min=1"`
	Subdomains []Subdomain `json:"subdomains" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type UpdateDomainForm struct {
	Name       string `json:"name" binding:"omitempty,min=1"`
	NameServer string `json:"nameServer" binding:"omitempty,min=1"`
	NSIp       string `json:"nsIp" binding:"omitempty,ip"`
	Ttl        string `json:"ttl" binding:"omitempty,min=1"`
}
