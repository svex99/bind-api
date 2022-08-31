package models

import (
	"fmt"
	"os"

	"github.com/svex99/bind-api/pkg/file"
	"github.com/svex99/bind-api/pkg/setting"
	"github.com/svex99/bind-api/services"
	"gorm.io/gorm"
)

type Domain struct {
	Id         uint        `json:"id" gorm:"not null;primaryKey"`
	Name       string      `json:"name" binding:"min=1" gorm:"not null;unique"`
	NameServer string      `json:"nameServer" binding:"min=1" gorm:"not null"`
	NSIp       string      `json:"nsIp" binding:"ip" gorm:"not null"`
	Ttl        string      `json:"ttl" binding:"min=1" gorm:"not null"`
	Subdomains []Subdomain `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	SOARecord  *SOARecord  `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	NSRecords  []NSRecord  `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ARecords   []ARecord   `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type UpdateDomainForm struct {
	Name       string `json:"name" binding:"omitempty,min=1"`
	NameServer string `json:"nameServer" binding:"omitempty,min=1"`
	NSIp       string `json:"nsIp" binding:"omitempty,ip"`
	Ttl        string `json:"ttl" binding:"omitempty,min=1"`
}

// Returns the path to the file that contains the records for the domain.
func (d *Domain) getFilePath() string {
	return setting.Bind.RecordsPath + "db." + d.Name
}

// Returns the ORIGIN string for the domain.
func (d *Domain) getOriginString() string {
	return fmt.Sprintf("$ORIGIN %s.\n", d.Name)
}

// Returns the TTL string for the domain.
func (d *Domain) getTtlString() string {
	return fmt.Sprintf("TTL %s\n", d.Ttl)
}

// Returns the zone string for the domain.
func (d *Domain) getZoneString() string {
	return fmt.Sprintf(
		"zone \"%s\" {\n\ttype master;\n\tfile \"/var/lib/bind/db.%s\";\n};\n",
		d.Name, d.Name,
	)
}

// Creates a new domain.
// This should create the necessary records and reload the bind configuration.
func (d *Domain) Create() error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// Create the domain in DB
		if err := tx.Create(d).Error; err != nil {
			return err
		}

		baseRecord := Record{
			Ttl:      d.Ttl,
			Class:    "IN",
			DomainId: d.Id,
		}

		// Create the SOA record for the domain in DB
		soaRecord := SOARecord{
			Record:     baseRecord,
			NameServer: d.NameServer,
			Admin:      setting.Bind.Admin,
			Refresh:    604800,
			Retry:      86400,
			Expire:     2419200,
			Minimum:    604800,
		}
		soaRecord.updateSerial()

		if err := tx.Create(&soaRecord).Error; err != nil {
			return err
		}

		// Create the NS record for the domain in DB
		nsRecord := NSRecord{
			Record:     baseRecord,
			NameServer: d.NameServer,
		}

		if err := tx.Create(&nsRecord).Error; err != nil {
			return err
		}

		// Create the A record for the name server in DB
		aRecord := ARecord{
			Record: baseRecord,
			Name:   d.NameServer,
			Ip:     d.NSIp,
		}

		if err := tx.Create(&aRecord).Error; err != nil {
			return err
		}

		// Set bind configuration files
		domainFile, err := os.Create(d.getFilePath())
		if err != nil {
			return err
		}

		defer domainFile.Close()

		if _, err := domainFile.WriteString(fmt.Sprintf(
			"%s%s%s%s%s", d.getOriginString(), d.getTtlString(), soaRecord.String(), nsRecord.String(), aRecord.String(),
		)); err != nil {
			return err
		}

		zoneFile, err := os.OpenFile(services.Bind.ZoneFilePath, os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			return err
		}

		defer zoneFile.Close()

		if _, err := zoneFile.WriteString(d.getZoneString()); err != nil {
			return err
		}

		// Reload the bind service
		if err := services.Bind.Reload(); err != nil {
			return err
		}

		return nil
	})
}

func (d *Domain) Update(form *UpdateDomainForm) error {
	if err := DB.Preload("SOARecord").First(d).Error; err != nil {
		return err
	}

	nsRecord := &NSRecord{}
	if err := DB.Where("domain_id = ? AND name_server = ?", d.Id, d.NameServer).First(nsRecord).Error; err != nil {
		return err
	}

	aRecord := &ARecord{}
	if err := DB.Where("domain_id = ? AND name = ? and ip = ?", d.Id, d.NameServer, d.NSIp).First(aRecord).Error; err != nil {
		return err
	}

	oldOriginString := d.getOriginString()
	oldTtlString := d.getTtlString()
	oldSOAString := d.SOARecord.String()
	oldNSString := nsRecord.String()
	oldAString := aRecord.String()

	if form.Name != "" {
		d.Name = form.Name
	}
	if form.NameServer != "" {
		d.NameServer = form.NameServer
	}
	if form.NSIp != "" {
		d.NSIp = form.NSIp
	}
	if form.Ttl != "" {
		d.Ttl = form.Ttl
	}

	return DB.Transaction(func(tx *gorm.DB) error {
		// Update domain in DB
		if err := tx.Save(d).Error; err != nil {
			return err
		}

		// Update SOA record in DB
		d.SOARecord.updateSerial()
		d.SOARecord.NameServer = d.NameServer

		if err := tx.Save(&d.SOARecord).Error; err != nil {
			return err
		}

		// Update NS record in DB
		nsRecord.NameServer = d.NameServer

		if err := tx.Save(nsRecord).Error; err != nil {
			return err
		}

		// Update A record in DB
		aRecord.Name = d.NameServer
		aRecord.Ip = d.NSIp

		if err := tx.Save(aRecord).Error; err != nil {
			return err
		}

		// Update bind configuration files
		if err := file.ReplaceContent(d.getFilePath(), oldOriginString, d.getOriginString(), false); err != nil {
			return err
		}

		if err := file.ReplaceContent(d.getFilePath(), oldTtlString, d.getTtlString(), false); err != nil {
			return err
		}

		if err := file.ReplaceContent(d.getFilePath(), oldSOAString, d.SOARecord.String(), true); err != nil {
			return err
		}

		if err := file.ReplaceContent(d.getFilePath(), oldNSString, nsRecord.String(), false); err != nil {
			return err
		}

		if err := file.ReplaceContent(d.getFilePath(), oldAString, aRecord.String(), false); err != nil {
			return err
		}

		// Reload bind service
		if err := services.Bind.Reload(); err != nil {
			return err
		}

		return nil
	})
}

func (d *Domain) Delete() error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// Load domain to have access to its Name in getFilePath()
		if err := tx.First(d).Error; err != nil {
			return err
		}

		// Remove the domain from DB
		// This should delete its related subdomains and records
		if err := tx.Delete(d).Error; err != nil {
			return err
		}

		// Remove the domain records file
		if err := os.Rename(d.getFilePath(), d.getFilePath()+".bak"); err != nil {
			return err
		}

		if err := file.ReplaceContent(services.Bind.ZoneFilePath, d.getZoneString(), "", true); err != nil {
			return err
		}

		// Reload the bind service
		if err := services.Bind.Reload(); err != nil {
			return err
		}

		return nil
	})
}
