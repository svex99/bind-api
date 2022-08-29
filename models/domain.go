package models

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

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

// temporal method to get the zone string for a domain
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

		origin := fmt.Sprintf("$ORIGIN %s.", d.Name)
		ttl := fmt.Sprintf("$TTL %s", d.Ttl)
		soa := fmt.Sprintf(
			"@ %s SOA %s %s ( %d %d %d %d %d )",
			soaRecord.Class, soaRecord.NameServer, soaRecord.Admin,
			soaRecord.Serial, soaRecord.Refresh, soaRecord.Retry, soaRecord.Expire, soaRecord.Minimum,
		)
		ns := fmt.Sprintf("@ %s NS %s", nsRecord.Class, nsRecord.NameServer)
		a := fmt.Sprintf("%s %s A %s", aRecord.Name, aRecord.Class, aRecord.Ip)

		if _, err := domainFile.WriteString(
			fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n", origin, ttl, soa, ns, a),
		); err != nil {
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
	oldDom := d

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

		// Get and update SOA record in DB
		if err := tx.Preload("SOARecord").Find(d).Error; err != nil {
			return err
		}

		d.SOARecord.updateSerial()
		d.SOARecord.NameServer = d.NameServer

		if err := tx.Save(&d.SOARecord).Error; err != nil {
			return err
		}

		// Get and update NS record in DB
		nsRecord := &NSRecord{}
		if err := tx.Where("domain_id = ? AND name_server = ?", d.Id, oldDom.NameServer).Find(nsRecord).Error; err != nil {
			return err
		}

		nsRecord.NameServer = d.NameServer

		if err := tx.Save(nsRecord).Error; err != nil {
			return err
		}

		// Get and update name server A record in DB
		aRecord := &ARecord{}
		if err := tx.Where("domain_id = ? AND name = ? and ip = ?", d.Id, oldDom.Name, oldDom.NSIp).Find(aRecord).Error; err != nil {
			return err
		}

		aRecord.Name = d.NameServer
		aRecord.Ip = d.NSIp

		if err := tx.Save(aRecord).Error; err != nil {
			return err
		}

		// TODO: Update bind configuration files
		// domainFile, err := os.OpenFile(setting.Bind.RecordsPath+"db."+d.Name, os.O_RDWR|os.O_APPEND, 0644)
		// if err != nil {
		// 	return err
		// }

		// defer domainFile.Close()

		// zoneFile, err := os.OpenFile(setting.Bind.ConfPath+"db."+oldDom.Name, os.O_RDWR|os.O_APPEND, 0644)
		// if err != nil {
		// 	return err
		// }

		// defer zoneFile.Close()

		// Reload bind service
		// if err := services.Bind.Reload(); err != nil {
		// 	return err
		// }

		return nil
	})
}

func (d *Domain) Delete() error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// Load domain to have access to its Name in getFilePath()
		if err := tx.Find(d).Error; err != nil {
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

		// Load the zones file content and update its data
		zoneContent, err := ioutil.ReadFile(services.Bind.ZoneFilePath)
		if err != nil {
			return err
		}

		updatedContent := strings.Replace(string(zoneContent), d.getZoneString(), "", 1)

		// Rename old zone file for backup
		if err := os.Rename(services.Bind.ZoneFilePath, services.Bind.ZoneFilePath+".bak"); err != nil {
			return err
		}

		// Create updated zone file
		if err := os.WriteFile(services.Bind.ZoneFilePath, []byte(updatedContent), 0666); err != nil {
			return err
		}

		// Reload the bind service
		if err := services.Bind.Reload(); err != nil {
			return err
		}

		return nil
	})
}
