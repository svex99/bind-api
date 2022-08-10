package models

type Subdomain struct {
	Id       uint64 `gorm:"not null;primaryKey"`
	Name     string `json:"name" binding:"required,min=1" gorm:"not null;unique"`
	DomainId uint64 `uri:"domain_id" gorm:"not null"`
	Ip       string `json:"ip" binding:"required,ip" gorm:"not null"`
}

type SubdomainInfo struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}

type UpdateSubdomainForm struct {
	Name string `json:"name" binding:"omitempty,min=1"`
	Ip   string `json:"ip" binding:"omitempty,ip"`
}
