package parser_test

// import (
// 	"os"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/svex99/bind-api/pkg/parser/domainParser"
// )

// func TestBuildUpdates(t *testing.T) {
// 	file, err := os.Open("testdata/db.test.com")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	parsedConf, err := domainParser.Parser.Parse("", file)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// TODO: Add test for SOA, MX, TXT, PTR
// 	nsRec := parsedConf.Records[0].(domainParser.NSRecord)
// 	aRec := parsedConf.Records[1].(domainParser.ARecord)

// 	data := []interface{}{
// 		domainParser.NSRecord{
// 			NameServer: domainParser.Name{Value: "new-name-server"},
// 			Id:         "NS1",
// 		},
// 		domainParser.ARecord{
// 			Name: domainParser.Name{Value: "new-name"},
// 			Ip:   domainParser.Ipv4{Value: "15.15.15.15"},
// 			Id:   "A1",
// 		},
// 	}

// 	builtUpdates, err := domainParser.BuildUpdates(parsedConf, data)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	expectedUpdates := []*domainParser.Update{
// 		{Offset: nsRec.NameServer.Pos.Offset, Len: 3, Content: "new-name-server"},
// 		{Offset: aRec.Name.Pos.Offset, Len: 3, Content: "new-name"},
// 		{Offset: aRec.Ip.Pos.Offset, Len: 11, Content: "15.15.15.15"},
// 	}

// 	assert.Equal(t, expectedUpdates, builtUpdates)

// 	t.Logf("->%s<-", nsRec.Pos.GoString())
// 	t.Logf("->%s<-", aRec.Pos.GoString())
// 	t.Logf("->%s<-", parsedConf.SOARecord.Pos.GoString())
// 	t.Log(nsRec.GetOffsetAndLen())
// 	t.Log(parsedConf)
// }
