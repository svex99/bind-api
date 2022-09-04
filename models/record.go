package models

import (
	"fmt"
	"time"
)

type Record struct {
	Id       uint   `json:"id" gorm:"primaryKey"`
	DomainId uint   `json:"domainId" binding:"min=1" gorm:"not null"`
	Domain   Domain `json:"-" binding:"-"`
}

type SOARecord struct {
	Record
	NameServer string `json:"nameServer" binding:"min=1"`
	Admin      string `json:"admin" binding:"min=1"`
	Serial     uint   `json:"serial"`
	Refresh    uint   `json:"refresh" binding:"gt=0"`
	Retry      uint   `json:"retry" binding:"gt=0"`
	Expire     uint   `json:"expire" binding:"gt=0"`
	Minimum    uint   `json:"minimum" binding:"gt=0"`
}

type NSRecord struct {
	Record
	NameServer string `json:"nameServer" binding:"min=1"`
}

type ARecord struct {
	Record
	Name string `json:"name" binding:"min=1"`
	Ip   string `json:"ip" binding:"ipv4"`
}

type MXRecord struct {
	Record
	Priority    uint   `json:"priority" binding:"gt=0"`
	EmailServer string `json:"emailServer" binding:"min=1"`
}

// Generates a new serial for the SOA record.
// Generated serials follows the format YYYYMMDDNN where NN is a two digits identifier.
func (soa *SOARecord) updateSerial() {
	now := time.Now().UTC()
	newSerial := uint(now.Year()*1_000_000 + int(now.Month())*10_000 + now.Day()*100)

	if soa.Serial >= newSerial {
		soa.Serial = soa.Serial + 1
	} else {
		soa.Serial = newSerial
	}
}

func (soa *SOARecord) String() string {
	return fmt.Sprintf(
		"@ IN SOA %s %s ( %d %d %d %d %d )\n",
		soa.NameServer, soa.Admin,
		soa.Serial, soa.Refresh, soa.Retry, soa.Expire, soa.Minimum,
	)
}

func (ns *NSRecord) String() string {
	return fmt.Sprintf("@ IN NS %s\n", ns.NameServer)
}

func (a *ARecord) String() string {
	return fmt.Sprintf("%s IN A %s\n", a.Name, a.Ip)
}

func (mx *MXRecord) String() string {
	return fmt.Sprintf("@ IN MX %d %s\n", mx.Priority, mx.EmailServer)
}
