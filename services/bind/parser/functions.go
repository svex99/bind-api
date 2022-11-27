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
func (zc *ZoneConf) WriteToDisk(filename string) (func(), error) {
	zc.UpdateSerial()

	content := []string{
		fmt.Sprintf("$ORIGIN %s.", zc.Origin),
		fmt.Sprintf("$TTL %s", zc.Ttl),
		fmt.Sprintf(
			"@ IN SOA %s %s ( %d %d %d %d %d )\n",
			zc.SOARecord.NameServer, zc.SOARecord.Admin,
			zc.SOARecord.Serial, zc.SOARecord.Refresh, zc.SOARecord.Retry, zc.SOARecord.Expire, zc.SOARecord.Minimum,
		),
	}
	for _, record := range zc.Records {
		content = append(content, record.String())
	}

	// Create a backup of config if file exists
	rollback := file.MakeBackup(filename)

	if err := os.WriteFile(filename, []byte(strings.Join(content, "\n")), 0666); err != nil {
		return rollback, err
	}

	return rollback, nil
}

func (zc *ZoneConf) DeleteFromDisk(filename string) (func(), error) {
	// Create a backup of config if file exists
	rollback := file.MakeBackup(filename)

	// File is already deleted when made the backup, since it's renamed.

	return rollback, nil
}

func (zc *ZoneConf) GetRecordIndex(targetRecord Record) int {
	target := targetRecord.String()
	return zc.GetRecordStringIndex(target)
}

func (zc *ZoneConf) GetRecordStringIndex(target string) int {
	for i, record := range zc.Records {
		if record.String() == target {
			return i
		}
	}

	return -1
}

// Generates a new serial for the SOA record.
// Generated serials follows the format YYYYMMDDNN where NN is a two digits identifier.
func (zc *ZoneConf) UpdateSerial() {
	now := time.Now().UTC()
	newSerial := uint(now.Year()*1_000_000 + int(now.Month())*10_000 + now.Day()*100)

	if zc.SOARecord.Serial >= newSerial {
		zc.SOARecord.Serial = zc.SOARecord.Serial + 1
	} else {
		zc.SOARecord.Serial = newSerial
	}
}

func (zc *ZoneConf) AddRecord(record Record) error {
	index := zc.GetRecordIndex(record)
	if index != -1 {
		return fmt.Errorf("record '%s' exists already", record.String())
	}

	zc.Records = append(zc.Records, record)

	return nil
}

func (zc *ZoneConf) UpdateRecord(target string, record Record) error {
	index := zc.GetRecordStringIndex(target)
	if index == -1 {
		return fmt.Errorf("target '%s' record does not exist", target)
	}

	zc.Records[index] = record

	if rollback, err := zc.WriteToDisk(zc.GetFilename()); err != nil {
		rollback()
		return err
	}

	return nil
}

func (zc *ZoneConf) DeleteRecord(record Record) error {
	index := zc.GetRecordIndex(record)
	if index == -1 {
		return fmt.Errorf("record '%s' does not exist", record.String())
	}

	zc.Records = append(zc.Records[:index], zc.Records[index+1:]...)

	return nil
}
