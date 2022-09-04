package models

import (
	"fmt"

	"github.com/svex99/bind-api/pkg/file"
	"github.com/svex99/bind-api/services"
	"gorm.io/gorm"
)

type TXTRecord struct {
	Record
	Value string `json:"value" binding:"min=1"`
}

type UpdateTXTRecordForm struct {
	Value string `json:"value" binding:"omitempty,min=1"`
}

func (txt *TXTRecord) String() string {
	return fmt.Sprintf("@ IN TXT %s\n", txt.Value)
}

func (txt *TXTRecord) Create() error {

	domain := &Domain{Id: txt.DomainId}

	if err := DB.Preload("SOARecord").First(domain).Error; err != nil {
		return err
	}

	soaRecord := domain.SOARecord

	oldSOAString := soaRecord.String()

	return DB.Transaction(func(tx *gorm.DB) error {
		// Save the new TXT record in DB
		if err := tx.Create(txt).Error; err != nil {
			return err
		}

		// Update SOA serial and save it in DB
		soaRecord.updateSerial()

		if err := tx.Save(soaRecord).Error; err != nil {
			return err
		}

		// Update bind configuration files
		if err := file.ReplaceContent(domain.getFilePath(), oldSOAString, soaRecord.String(), true); err != nil {
			return err
		}

		if err := file.AddContent(domain.getFilePath(), txt.String()); err != nil {
			return err
		}

		// Reload bind with new configuration
		if err := services.Bind.Reload(); err != nil {
			return err
		}

		return nil
	})
}

func (txt *TXTRecord) Update(form *UpdateTXTRecordForm) error {
	domain := &Domain{Id: txt.DomainId}

	soaRecord, err := domain.GetWithSOA()
	if err != nil {
		return err
	}

	if err := DB.Debug().Where(txt).First(txt).Error; err != nil {
		return err
	}

	oldSOAString := soaRecord.String()
	oldTXTString := txt.String()

	if form.Value != "" {
		txt.Value = form.Value
	}

	return DB.Transaction(func(tx *gorm.DB) error {
		// Update TXT record on DB
		if err := tx.Save(txt).Error; err != nil {
			return err
		}

		// Update SOA serial on DB
		soaRecord.updateSerial()

		if err := tx.Save(soaRecord).Error; err != nil {
			return err
		}

		// Update bind configuration files
		if err := file.ReplaceContent(domain.getFilePath(), oldSOAString, soaRecord.String(), true); err != nil {
			return err
		}

		if err := file.ReplaceContent(domain.getFilePath(), oldTXTString, txt.String(), true); err != nil {
			return err
		}

		// Reload bind service with new configuration
		if err := services.Bind.Reload(); err != nil {
			return err
		}

		return nil
	})
}

func (txt *TXTRecord) Delete() error {
	domain := Domain{Id: txt.DomainId}

	soaRecord, err := domain.GetWithSOA()
	if err != nil {
		return err
	}

	if err := DB.Where(txt).First(txt).Error; err != nil {
		return err
	}

	oldSOAString := soaRecord.String()

	return DB.Transaction(func(tx *gorm.DB) error {
		// Delete TXT record on DB
		if err := tx.Delete(txt).Error; err != nil {
			return err
		}

		// Update SOA serial on DB
		soaRecord.updateSerial()

		if err := tx.Save(soaRecord).Error; err != nil {
			return err
		}

		// Update bind configuration files
		if err := file.ReplaceContent(domain.getFilePath(), oldSOAString, soaRecord.String(), true); err != nil {
			return err
		}

		if err := file.ReplaceContent(domain.getFilePath(), txt.String(), "", false); err != nil {
			return err
		}

		// Reload bind service with new configuration
		if err := services.Bind.Reload(); err != nil {
			return err
		}

		return nil
	})
}
