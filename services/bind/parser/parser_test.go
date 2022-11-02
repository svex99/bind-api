package parser_test

// import (
// 	"os"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/svex99/bind-api/pkg/parser/domainParser"
// )

// func TestParser(t *testing.T) {
// 	file, err := os.Open("testdata/db.test.com")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	t.Log(domainParser.Parser.String())

// 	conf, err := domainParser.Parser.Parse("", file)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	assert.Equal(t, "test.com", conf.Origin.Value)
// 	assert.Equal(t, "12d", conf.Ttl.Value)

// 	soa := conf.SOARecord
// 	assert.Equal(t, "ns1", soa.NameServer.Value)
// 	assert.Equal(t, "admin", soa.Admin.Value)
// 	assert.Equal(t, "1", soa.Serial.Value)
// 	assert.Equal(t, "23", soa.Refresh.Value)
// 	assert.Equal(t, "456", soa.Retry.Value)
// 	assert.Equal(t, "78", soa.Expire.Value)
// 	assert.Equal(t, "90", soa.Minimum.Value)

// 	records := conf.Records

// 	assert.Len(t, records, 7)

// 	nsRec, ok := records[0].(domainParser.NSRecord)
// 	if !ok {
// 		t.Fatal("Records[0] must be a record of type NS")
// 	}
// 	assert.Equal(t, "ns1", nsRec.NameServer.Value)
// 	assert.Equal(t, "NS1", nsRec.Id)

// 	aRec, ok := records[1].(domainParser.ARecord)
// 	if !ok {
// 		t.Fatal("Records[1] must be a record of type A")
// 	}
// 	assert.Equal(t, "ns1", aRec.Name.Value)
// 	assert.Equal(t, "10.10.10.10", aRec.Ip.Value)
// 	assert.Equal(t, "A1", aRec.Id)

// 	mxRec, ok := records[2].(domainParser.MXRecord)
// 	if !ok {
// 		t.Fatal("Records[2] must be a record of type MX")
// 	}
// 	assert.Equal(t, "100", mxRec.Priority.Value)
// 	assert.Equal(t, "email", mxRec.EmailServer.Value)
// 	assert.Equal(t, "MX1", mxRec.Id)

// 	txtRec, ok := records[3].(domainParser.TXTRecord)
// 	if !ok {
// 		t.Fatal("Records[3] must be a record of type TXT")
// 	}
// 	assert.Equal(t, "var=hello world1", txtRec.Value.Value)
// 	assert.Equal(t, "TXT1", txtRec.Id)
// }
