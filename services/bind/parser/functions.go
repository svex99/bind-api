package parser

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/svex99/bind-api/pkg/file"
)

// Writes the zone configuration to a plain text file.
// Returns a function that rollbacks the process.
func (dc *DomainConf) WriteToDisk(filename string) (func(), error) {
	dc.UpdateSerial()

	content := []string{
		fmt.Sprintf("$ORIGIN %s.", dc.Origin),
		fmt.Sprintf("$TTL %s", dc.Ttl),
		fmt.Sprintf(
			"@ IN SOA %s %s ( %d %d %d %d %d )",
			dc.SOARecord.NameServer, dc.SOARecord.Admin,
			dc.SOARecord.Serial, dc.SOARecord.Refresh, dc.SOARecord.Retry, dc.SOARecord.Expire, dc.SOARecord.Minimum,
		),
	}
	for _, record := range dc.Records {
		content = append(content, record.String())
	}

	// Create a backup of config if file exists
	rollback := file.MakeBackup(filename)

	if err := os.WriteFile(filename, []byte(strings.Join(content, "\n")), 0666); err != nil {
		return rollback, err
	}

	return rollback, nil
}

func (dc *DomainConf) DeleteFromDisk(filename string) (func(), error) {
	// Create a backup of config if file exists
	rollback := file.MakeBackup(filename)

	// File is already deleted when made the backup, since it's renamed.

	return rollback, nil
}

func (dc *DomainConf) GetRecordIndex(hash uint) (int, error) {
	for i, record := range dc.Records {
		if record.GetHash() == hash {
			return i, nil
		}
	}

	return -1, fmt.Errorf("no record found for hash '%d'", hash)
}

// Generates a new serial for the SOA record.
// Generated serials follows the format YYYYMMDDNN where NN is a two digits identifier.
func (dc *DomainConf) UpdateSerial() {
	now := time.Now().UTC()
	newSerial := uint(now.Year()*1_000_000 + int(now.Month())*10_000 + now.Day()*100)

	if dc.SOARecord.Serial >= newSerial {
		dc.SOARecord.Serial = dc.SOARecord.Serial + 1
	} else {
		dc.SOARecord.Serial = newSerial
	}
}

func (dc *DomainConf) AddRecord(newRecord Record) error {
	recordHash := newRecord.GetHash()

	_, err := dc.GetRecordIndex(recordHash)
	if err == nil {
		return fmt.Errorf("record with hash '%d' exists already", recordHash)
	}

	dc.Records = append(dc.Records, newRecord)

	if rollback, err := dc.WriteToDisk(dc.GetFilename()); err != nil {
		rollback()
		return err
	}

	return nil
}

func (dc *DomainConf) UpdateRecord(targetHash uint, updateRecord Record) error {
	targetIndex, err := dc.GetRecordIndex(targetHash)
	if err != nil {
		return err
	}

	dc.Records[targetIndex] = updateRecord

	if rollback, err := dc.WriteToDisk(dc.GetFilename()); err != nil {
		rollback()
		return err
	}

	return nil
}

func (dc *DomainConf) DeleteRecord(targetHash uint) error {
	targetIndex, err := dc.GetRecordIndex(targetHash)
	if err != nil {
		return err
	}

	dc.Records = append(dc.Records[:targetIndex], dc.Records[targetIndex+1:]...)

	if rollback, err := dc.WriteToDisk(dc.GetFilename()); err != nil {
		rollback()
		return err
	}

	return nil
}
