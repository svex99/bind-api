package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/svex99/bind-api/services/bind"
	"github.com/svex99/bind-api/services/bind/parser"
)

func getTargetFromRequest(c *gin.Context) (*parser.DomainConf, parser.Record, error) {
	origin := c.Param("origin")
	recordType := strings.ToUpper(c.Param("type"))

	var data map[string]string

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, nil, err
	}

	dConf, ok := bind.Service.Domains[origin]
	if !ok {
		err := fmt.Errorf("domain %s was not found", origin)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return nil, nil, err
	}

	var record parser.Record

	switch recordType {
	case "NS":
		record = parser.NSRecord{Type: recordType, NameServer: string(data["nameServer"])}
	case "A":
		record = parser.ARecord{Type: recordType, Name: data["name"], Ip: data["ip"]}
	case "MX":
		priority, err := strconv.ParseUint(data["priority"], 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return nil, nil, err
		}
		record = parser.MXRecord{Type: recordType, Priority: uint(priority), EmailServer: data["emailServer"]}
	case "TXT":
		record = parser.TXTRecord{Type: recordType, Value: data["value"]}
	case "CNAME":
		record = parser.CNAMERecord{Type: recordType, SrcName: data["srcName"], DstName: data["dstName"]}
	}

	return dConf, record, nil
}

func PostRecord(c *gin.Context) {
	bind.Service.Mutex.Lock()
	defer bind.Service.Mutex.Unlock()

	dConfPtr, record, err := getTargetFromRequest(c)
	if err != nil {
		return
	}

	dConf := *dConfPtr

	if err := dConf.AddRecord(record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bind.Service.Domains[c.Param("origin")] = &dConf

	c.JSON(http.StatusCreated, gin.H{"record": record})
}

func PatchRecord(c *gin.Context) {
	bind.Service.Mutex.Lock()
	defer bind.Service.Mutex.Unlock()

	hash, err := strconv.ParseUint(c.Param("hash"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("invalid hash %s", c.Param("hash"))})
	}

	dConfPtr, record, err := getTargetFromRequest(c)
	if err != nil {
		return
	}

	dConf := *dConfPtr

	if err := dConf.UpdateRecord(uint(hash), record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bind.Service.Domains[c.Param("origin")] = &dConf

	c.JSON(http.StatusOK, gin.H{"record": record})
}

func DeleteRecord(c *gin.Context) {
	bind.Service.Mutex.Lock()
	defer bind.Service.Mutex.Unlock()

	dConfPtr, record, err := getTargetFromRequest(c)
	if err != nil {
		return
	}

	dConf := *dConfPtr

	if err := dConf.DeleteRecord(record.GetHash()); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	bind.Service.Domains[c.Param("origin")] = &dConf

	c.JSON(http.StatusOK, gin.H{"message": "record deleted"})
}
