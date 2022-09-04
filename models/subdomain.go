package models

import (
	"github.com/svex99/bind-api/pkg/file"
	"github.com/svex99/bind-api/services"
	"gorm.io/gorm"
)

type Subdomain struct {
	Id   uint   `json:"id"`
	Name string `json:"name" binding:"required,min=1" gorm:"not null;unique"`
	Ip   string `json:"ip" binding:"ipv4" gorm:"not null"`
}

type UpdateSubdomainForm struct {
	Name string `json:"name" binding:"omitempty,min=1"`
	Ip   string `json:"ip" binding:"omitempty,ip"`
}

func (s *Subdomain) Create(domainId uint) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// Get the Domain and its SOA record from DB
		domain := &Domain{Id: domainId}

		if err := DB.Preload("SOARecord").First(domain).Error; err != nil {
			return err
		}

		soaRecord := domain.SOARecord

		// Create the A record in DB
		aRecord := ARecord{
			Record: Record{
				DomainId: domain.Id,
			},
			Name: s.Name,
			Ip:   s.Ip,
		}

		if err := tx.Create(&aRecord).Error; err != nil {
			return err
		}

		s.Id = aRecord.Id

		oldSOAString := soaRecord.String()

		// Update the SOA record serial and save on DB
		soaRecord.updateSerial()

		if err := tx.Save(soaRecord).Error; err != nil {
			return err
		}

		// Update bind configuration files
		if err := file.ReplaceContent(domain.getFilePath(), oldSOAString, soaRecord.String(), true); err != nil {
			return err
		}

		if err := file.AddContent(domain.getFilePath(), aRecord.String()); err != nil {
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
	aRecord := &ARecord{Record: Record{Id: s.Id, DomainId: domainId}}

	if err := DB.Preload("Domain").First(aRecord, aRecord).Error; err != nil {
		return err
	}

	domain := &aRecord.Domain

	if err := DB.Preload("SOARecord").First(domain).Error; err != nil {
		return err
	}

	soaRecord := domain.SOARecord

	oldSOAString := soaRecord.String()
	oldAString := aRecord.String()

	if form.Name != "" {
		aRecord.Name = form.Name
	}
	if form.Ip != "" {
		aRecord.Ip = form.Ip
	}

	s.Name = aRecord.Name
	s.Ip = aRecord.Ip

	return DB.Transaction(func(tx *gorm.DB) error {
		// Update A record in DB
		if err := tx.Save(aRecord).Error; err != nil {
			return err
		}

		// Update SOA record in DB
		soaRecord.updateSerial()

		if err := tx.Save(soaRecord).Error; err != nil {
			return err
		}

		// Update bind configuration files
		if err := file.ReplaceContent(domain.getFilePath(), oldSOAString, soaRecord.String(), true); err != nil {
			return err
		}

		if err := file.ReplaceContent(domain.getFilePath(), oldAString, aRecord.String(), false); err != nil {
			return err
		}

		// Reload bind service with new configuration
		if err := services.Bind.Reload(); err != nil {
			return err
		}

		return nil
	})
}

func (s *Subdomain) Delete(domainId uint) error {
	aRecord := &ARecord{Record: Record{Id: s.Id, DomainId: domainId}}

	if err := DB.Preload("Domain").First(aRecord, aRecord).Error; err != nil {
		return err
	}

	domain := &aRecord.Domain

	if err := DB.Preload("SOARecord").First(domain).Error; err != nil {
		return err
	}

	soaRecord := domain.SOARecord

	oldSOAString := soaRecord.String()
	oldAString := aRecord.String()

	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(aRecord).Error; err != nil {
			return err
		}

		// Update SOA serial in DB
		soaRecord.updateSerial()

		if err := tx.Save(soaRecord).Error; err != nil {
			return err
		}

		// Update bind configuration files
		if err := file.ReplaceContent(domain.getFilePath(), oldSOAString, soaRecord.String(), true); err != nil {
			return err
		}

		if err := file.ReplaceContent(domain.getFilePath(), oldAString, "", false); err != nil {
			return err
		}

		// Reload bind service with new configuration
		if err := services.Bind.Reload(); err != nil {
			return err
		}

		return nil
	})
}
