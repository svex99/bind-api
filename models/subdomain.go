package models

import (
	"github.com/svex99/bind-api/pkg/file"
	"github.com/svex99/bind-api/services"
	"gorm.io/gorm"
)

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
				Ttl:      s.Domain.Ttl,
				Class:    "IN",
				DomainId: s.Domain.Id,
			},
			Name: s.Name,
			Ip:   s.Ip,
		}

		if err := tx.Create(&aRecord).Error; err != nil {
			return err
		}

		// Set bind configuration files
		if err := file.AddContent(s.Domain.getFilePath(), aRecord.String()); err != nil {
			return err
		}

		// Reload bind service
		if err := services.Bind.Reload(); err != nil {
			return err
		}

		return nil
	})
}

func (s *Subdomain) Update(domainId uint, form *UpdateSubdomainForm) error {
	if err := DB.Preload("Domain").First(s, "domain_id = ?", domainId).Error; err != nil {
		return err
	}

	// Get SOA record of the domain, needed to update serial
	if err := DB.Preload("SOARecord").First(&s.Domain).Error; err != nil {
		return err
	}

	// Get A record of the subdomain
	aRecord := &ARecord{}
	if err := DB.First(aRecord, "domain_id = ?", domainId, "name = ?", s.Name, "ip = ?", s.Ip).Error; err != nil {
		return err
	}

	oldSOAString := s.Domain.SOARecord.String()
	oldAString := aRecord.String()

	if form.Name != "" {
		s.Name = form.Name
	}
	if form.Ip != "" {
		s.Ip = form.Ip
	}

	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(s).Error; err != nil {
			return err
		}

		// Update A record in DB
		aRecord.Name = s.Name
		aRecord.Ip = s.Ip

		if err := tx.Save(aRecord).Error; err != nil {
			return err
		}

		// Update SOA record in DB
		s.Domain.SOARecord.updateSerial()

		if err := tx.Save(&s.Domain.SOARecord).Error; err != nil {
			return err
		}

		// Update bind configuration files
		if err := file.ReplaceContent(s.Domain.getFilePath(), oldSOAString, s.Domain.SOARecord.String(), true); err != nil {
			return err
		}

		if err := file.ReplaceContent(s.Domain.getFilePath(), oldAString, aRecord.String(), false); err != nil {
			return err
		}

		// Reload bind service with new configuration
		if err := services.Bind.Reload(); err != nil {
			return err
		}

		return nil
	})
}
