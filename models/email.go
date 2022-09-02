package models

import (
	"github.com/svex99/bind-api/pkg/file"
	"github.com/svex99/bind-api/services"
	"gorm.io/gorm"
)

type Email struct {
	Id       uint   `json:"id" gorm:"not null;primaryKey"`
	Priority uint   `json:"priority" gorm:"not null" binding:"gt=0"`
	Name     string `json:"name" gorm:"not null" binding:"min=1"`
	Ip       string `json:"ip" gorm:"not null" binding:"ipv4"`
}

type UpdateEmailForm struct {
	Priority uint   `json:"priority" binding:"omitempty,gt=0"`
	Name     string `json:"name" binding:"omitempty,min=1"`
	Ip       string `json:"ip" binding:"omitempty,ipv4"`
}

func (e *Email) Create(domainId uint) error {
	// Preload its Domain and the SOARecord
	domain := &Domain{Id: domainId}

	if err := DB.Preload("SOARecord").First(domain).Error; err != nil {
		return err
	}

	oldSOAString := domain.SOARecord.String()

	return DB.Transaction(func(tx *gorm.DB) error {
		// Create the Email in DB
		// if err := tx.Create(e).Error; err != nil {
		// 	return err
		// }

		baseRecord := Record{
			Ttl:      domain.Ttl,
			Class:    "IN",
			DomainId: domain.Id,
		}

		// Create the MX record in DB
		mxRecord := &MXRecord{
			Record:      baseRecord,
			Priority:    e.Priority,
			EmailServer: e.Name,
		}

		if err := tx.Create(mxRecord).Error; err != nil {
			return err
		}

		e.Id = mxRecord.Id

		// TODO: Allow empty IP, that means to not add an A record
		// Create the A record in DB
		aRecord := &ARecord{
			Record: baseRecord,
			Name:   e.Name,
			Ip:     e.Ip,
		}

		if err := tx.Create(aRecord).Error; err != nil {
			return err
		}

		// Update bind configuration files
		if err := file.ReplaceContent(domain.getFilePath(), oldSOAString, domain.SOARecord.String(), true); err != nil {
			return err
		}

		if err := file.AddContent(domain.getFilePath(), mxRecord.String()); err != nil {
			return err
		}

		if err := file.AddContent(domain.getFilePath(), aRecord.String()); err != nil {
			return err
		}

		// Reload bind service with new configuration
		if err := services.Bind.Reload(); err != nil {
			return err
		}

		return nil
	})
}

func (e *Email) Update(domainId uint, form *UpdateEmailForm) error {
	// Get the MX record, A record, Domain and SOA record from DB
	mxRecord := &MXRecord{Record: Record{Id: e.Id}}

	if err := DB.Preload("Domain").First(mxRecord).Error; err != nil {
		return err
	}

	domain := &mxRecord.Domain

	if err := DB.Preload("SOARecord").First(&mxRecord.Domain).Error; err != nil {
		return err
	}

	soaRecord := domain.SOARecord

	// TODO: Can exist multiple A records for same MX record
	aRecord := &ARecord{}

	if err := DB.First(
		aRecord,
		&ARecord{
			Record: Record{
				DomainId: domain.Id,
			},
			Name: mxRecord.EmailServer,
		},
	).Error; err != nil {
		return err
	}

	oldSOAString := soaRecord.String()
	oldMXString := mxRecord.String()
	oldAString := aRecord.String()

	if form.Priority != 0 {
		mxRecord.Priority = form.Priority
		e.Priority = mxRecord.Priority
	}
	if form.Name != "" {
		mxRecord.EmailServer = form.Name
		aRecord.Name = form.Name
		e.Name = mxRecord.EmailServer
	}
	if form.Ip != "" {
		aRecord.Ip = form.Ip
		e.Ip = aRecord.Ip
	}

	return DB.Transaction(func(tx *gorm.DB) error {
		// Update SOA record serial and save in DB
		soaRecord.updateSerial()

		if err := tx.Save(soaRecord).Error; err != nil {
			return err
		}

		// Update MX record in DB
		if err := tx.Save(mxRecord).Error; err != nil {
			return err
		}

		// Update A record in DB
		if err := tx.Save(aRecord).Error; err != nil {
			return err
		}

		// Update bind configuration
		if err := file.ReplaceContent(domain.getFilePath(), oldSOAString, soaRecord.String(), true); err != nil {
			return err
		}

		if err := file.ReplaceContent(domain.getFilePath(), oldMXString, mxRecord.String(), false); err != nil {
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
