package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/svex99/bind-api/services/bind"
	"github.com/svex99/bind-api/services/bind/parser"
)

func getRecord(c *gin.Context) (parser.Record, error) {
	var data map[string]any

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, err
	}

	typeField, ok := data["type"]
	if !ok {
		err := errors.New("missing field 'type'")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, err
	}

	rType := strings.ToUpper(typeField.(string))

	var record parser.Record

	switch rType {
	case "NS":
		record = parser.NSRecord{Type: rType, NameServer: string(data["nameServer"].(string))}
	case "A":
		record = parser.ARecord{Type: rType, Name: data["name"].(string), Ip: data["ip"].(string)}
	case "MX":
		var priority uint

		priority, ok = data["priority"].(uint)
		if !ok {
			parsedPriority, err := strconv.ParseUint(data["priority"].(string), 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return nil, err
			}

			priority = uint(parsedPriority)
		}

		record = parser.MXRecord{Type: rType, Priority: priority, EmailServer: data["emailServer"].(string)}
	case "TXT":
		record = parser.TXTRecord{Type: rType, Value: data["value"].(string)}
	case "CNAME":
		record = parser.CNAMERecord{Type: rType, SrcName: data["srcName"].(string), DstName: data["dstName"].(string)}
	default:
		err := fmt.Errorf("field 'type' cannot be empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, err
	}

	return record, nil
}

func PostRecord(c *gin.Context) {
	origin := c.Param("origin")

	record, err := getRecord(c)
	if err != nil {
		return
	}

	if err := bind.Service.AddRecord(origin, record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, record)
}

func PatchRecord(c *gin.Context) {
	origin := c.Param("origin")
	target := c.Param("target") + "\n"

	log.Println(target)

	record, err := getRecord(c)
	if err != nil {
		return
	}

	if err := bind.Service.UpdateRecord(origin, target, record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, record)
}

func DeleteRecord(c *gin.Context) {
	origin := c.Param("origin")

	record, err := getRecord(c)
	if err != nil {
		return
	}

	if err := bind.Service.DeleteRecord(origin, record); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "record deleted"})
}
