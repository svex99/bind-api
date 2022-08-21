package models

import (
	"fmt"
	"strconv"
	"time"

	"github.com/svex99/bind-api/utils/setting"
	"gorm.io/gorm"
)

type Domain struct {
	Id         uint64      `json:"id" gorm:"not null;primaryKey"`
	Name       string      `json:"name" binding:"min=1" gorm:"not null;unique"`
	NameServer string      `json:"nameServer" binding:"min=1" gorm:"not null"`
	NSIp       string      `json:"nsIp" binding:"ip" gorm:"not null"`
	Ttl        string      `json:"ttl" binding:"min=1" gorm:"not null"`
	Subdomains []Subdomain `json:"subdomains" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type UpdateDomainForm struct {
	Name       string `json:"name" binding:"omitempty,min=1"`
	NameServer string `json:"nameServer" binding:"omitempty,min=1"`
	NSIp       string `json:"nsIp" binding:"omitempty,ip"`
	Ttl        string `json:"ttl" binding:"omitempty,min=1"`
}

// Creates a new domain.
// This should create the necessary records and reload the bind configuration.
func (d *Domain) Create() error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// Create the domain in DB
		if err := tx.Create(d).Error; err != nil {
			return err
		}

		// Create the SOA record for the domain in DB
		serial, err := newSerial()
		if err != nil {
			return err
		}

		soaRecord := SOARecord{
			Record: Record{
				Ttl:        d.Ttl,
				Class:      "IN",
				DomainName: d.Name,
			},
			NameServer: d.NameServer,
			Admin:      setting.App.BindAdmin,
			Serial:     serial,
			Refresh:    604800,
			Retry:      86400,
			Expire:     2419200,
			Minimum:    604800,
		}

		if err := tx.Create(&soaRecord).Error; err != nil {
			return err
		}

		// Create the NS record for the domain in DB
		nsRecord := NSRecord{
			Record: Record{
				Ttl:        d.Ttl,
				Class:      "IN",
				DomainName: d.Name,
			},
			NameServer: d.NameServer,
		}

		if err := tx.Create(&nsRecord).Error; err != nil {
			return err
		}

		// Create the A record for the name server in DB
		aRecord := ARecord{
			Record: Record{
				Ttl:        d.Ttl,
				Class:      "IN",
				DomainName: d.Name,
			},
			Name: d.NameServer,
			Ip:   d.NSIp,
		}

		if err := tx.Create(&aRecord).Error; err != nil {
			return err
		}

		// TODO: Set bind configuration files
		// TODO: Reload bind service

		return nil
	})
}

func newSerial() (uint, error) {
	now := time.Now().UTC()

	serialString := fmt.Sprintf(
		"%d%02d%02d%02d%02d%02d",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(),
	)

	serial, err := strconv.ParseUint(serialString, 10, 64)

	if err != nil {
		return 0, err
	}

	return uint(serial), nil
}
