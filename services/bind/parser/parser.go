package parser

import (
	"fmt"
	"hash/fnv"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/svex99/bind-api/pkg/setting"
)

var (
	Lexer = lexer.MustSimple([]lexer.SimpleRule{
		{Name: "Directive", Pattern: `\$(ORIGIN|TTL)`},
		{Name: "Keyword", Pattern: `@`},
		{Name: "RType", Pattern: `IN|SOA|NS|A|MX|TXT|PTR`},
		{Name: "Origin", Pattern: `[a-zA-Z][\w\-]*\.[a-zA-Z]+`},
		{Name: "Name", Pattern: `[a-zA-Z][\w\-]*`},
		{Name: "Ttl", Pattern: `\d+[hdw]`},
		{Name: "Ipv4", Pattern: `\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`},
		{Name: "Uint", Pattern: `\d+`},
		{Name: "String", Pattern: `"[^"\n]*"`},
		{Name: "Punct", Pattern: `[\.\(\)]`},
		{Name: "Comment", Pattern: `;[^\n]*\n+`},
		{Name: "Whitespace", Pattern: `[ \t\r]+`},
		{Name: "NewLine", Pattern: `[\n]+`},
	})
	Parser = participle.MustBuild[DomainConf](
		participle.Lexer(Lexer),
		participle.Union[Record](NSRecord{}, ARecord{}, MXRecord{}, TXTRecord{}, CNAMERecord{}),
		participle.Elide("Whitespace", "Comment"),
		participle.Unquote("String"),
		participle.UseLookahead(2),
	)
)

type DomainConf struct {
	Origin    string     `parser:"'$ORIGIN' @Origin '.' NewLine" json:"origin"`
	Ttl       string     `parser:"'$TTL' @Ttl NewLine" json:"ttl"`
	SOARecord *SOARecord `parser:"@@ NewLine" json:"soaRecord"`
	Records   []Record   `parser:"@@*" json:"records"`
}

// TODO: Do not return path from string concatenation
func (dc *DomainConf) GetFilename() string {
	return fmt.Sprintf("%sdb.%s", setting.Bind.RecordsPath, dc.Origin)
}

type Record interface {
	GetHash() uint
	String() string
}

type SOARecord struct {
	NameServer string `parser:"'@' 'IN' 'SOA' @Name" json:"nameServer"`
	Admin      string `parser:"@Name" json:"admin"`
	Serial     uint   `parser:"'(' @Uint" json:"serial"`
	Refresh    uint   `parser:"@Uint" json:"refresh"`
	Retry      uint   `parser:"@Uint" json:"retry"`
	Expire     uint   `parser:"@Uint" json:"expire"`
	Minimum    uint   `parser:"@Uint ')'" json:"minimum"`
}

type NSRecord struct {
	Type       string `parser:"'@' 'IN' @'NS'" json:"type"`
	NameServer string `parser:"@Name NewLine" json:"nameServer"`
}

func (ns NSRecord) GetHash() uint {
	h := fnv.New64a()
	h.Write([]byte(ns.String()))
	return uint(h.Sum64())
}

func (ns NSRecord) String() string {
	return fmt.Sprintf("@ IN NS %s\n", ns.NameServer)
}

type ARecord struct {
	Name string `parser:"@Name" json:"name"`
	Type string `parser:"'IN' @'A'" json:"type"`
	Ip   string `parser:"@Ipv4 NewLine" json:"ip"`
}

func (a ARecord) GetHash() uint {
	h := fnv.New64a()
	h.Write([]byte(a.String()))
	return uint(h.Sum64())
}

func (a ARecord) String() string {
	return fmt.Sprintf("%s IN A %s\n", a.Name, a.Ip)
}

type MXRecord struct {
	Type        string `parser:"'@' 'IN' @'MX'" json:"type"`
	Priority    uint   `parser:"@Uint" json:"priority"`
	EmailServer string `parser:"@Name NewLine" json:"emailServer"`
}

func (mx MXRecord) GetHash() uint {
	h := fnv.New64a()
	h.Write([]byte(mx.String()))
	return uint(h.Sum64())
}

func (mx MXRecord) String() string {
	return fmt.Sprintf("@ IN MX %d %s\n", mx.Priority, mx.EmailServer)
}

type TXTRecord struct {
	Type  string `parser:"'@' 'IN' @'TXT'" json:"type"`
	Value string `parser:"@String NewLine" json:"value"`
}

func (txt TXTRecord) GetHash() uint {
	h := fnv.New64a()
	h.Write([]byte(txt.String()))
	return uint(h.Sum64())
}

func (txt TXTRecord) String() string {
	return fmt.Sprintf("@ IN TXT %s", txt.Value)
}

type CNAMERecord struct {
	SrcName string `parser:"@Name 'IN'" json:"srcName"`
	Type    string `parser:"@'CNAME'" json:"type"`
	DstName string `parser:"@Name NewLine" json:"dstName"`
}

func (cname CNAMERecord) GetHash() uint {
	h := fnv.New64a()
	h.Write([]byte(cname.String()))
	return uint(h.Sum64())
}

func (cname CNAMERecord) String() string {
	return fmt.Sprintf("%s IN CNAME %s\n", cname.SrcName, cname.DstName)
}
