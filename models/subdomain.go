package models

import "gorm.io/gorm"

type Subdomain struct {
	Id       uint   `json:"id" gorm:"not null;primaryKey"`
	Name     string `json:"name" binding:"required,min=1" gorm:"not null;unique"`
	DomainId uint   `json:"domainId" gorm:"not null"`
	Domain   Domain `binding:"-"`
	Ip       string `json:"ip" binding:"ipv4" gorm:"not null"`
}

type UpdateSubdomainForm struct {
	Name string `json:"name" binding:"omitempty,min=1"`
	Ip   string `json:"ip" binding:"omitempty,ip"`
}

func (s *Subdomain) Create() error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// Create the subdomain in DB
		if err := tx.Create(s).Error; err != nil {
			return err
		}

		if err := tx.Preload("Domain").Find(s).Error; err != nil {
			return err
		}

		// Create the A record for the subdomain in DB
		aRecord := ARecord{
			Record: Record{
				Ttl:        s.Domain.Ttl,
				Class:      "IN",
				DomainName: s.Domain.Name,
			},
			Name: s.Name,
			Ip:   s.Ip,
		}

		if err := tx.Create(&aRecord).Error; err != nil {
			return err
		}

		// TODO: Set bind configuration files
		// TODO: Reload bind service

		return nil
	})
}
