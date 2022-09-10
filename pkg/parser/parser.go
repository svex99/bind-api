package parser

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

var (
	domainConfLexer = lexer.MustSimple([]lexer.SimpleRule{
		{Name: "Directive", Pattern: `\$(ORIGIN|TTL)`},
		{Name: "Keyword", Pattern: `@|IN|SOA|NS|A|MX|TXT|PTR`},
		{Name: "Domain", Pattern: `[a-zA-Z][\w\-]*\.[a-zA-Z]+`},
		{Name: "Name", Pattern: `[a-zA-Z][\w\-]*`},
		{Name: "Ttl", Pattern: `\d+[hdw]`},
		{Name: "Ipv4", Pattern: `\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`},
		{Name: "Uint", Pattern: `\d+`},
		{Name: "String", Pattern: `"[^"\n]*"`},
		{Name: "Punct", Pattern: `[\.\(\)]`},
		{Name: "Comment", Pattern: `;[^\n]*`},
		{Name: "Whitespace", Pattern: `[ \n\t\r]+`},
	})
	Parser = participle.MustBuild[DomainConf](
		participle.Lexer(domainConfLexer),
		participle.Union[Record](NSRecord{}, ARecord{}, MXRecord{}, TXTRecord{}),
		participle.Elide("Whitespace", "Comment"),
		participle.Unquote("String"),
		participle.UseLookahead(2),
	)
)

type DomainConf struct {
	Origin    Domain     `parser:"'$ORIGIN' @@ '.'"`
	Ttl       Ttl        `parser:"'$TTL' @@"`
	SOARecord *SOARecord `parser:"@@"`
	Records   []Record   `parser:"@@*"`
}

type Record interface{}

type SOARecord struct {
	Pos        lexer.Position
	NameServer Name `parser:"'@' 'IN' 'SOA' @@"`
	Admin      Name `parser:"@@"`
	Serial     Uint `parser:"'(' @@"`
	Refresh    Uint `parser:"@@"`
	Retry      Uint `parser:"@@"`
	Expire     Uint `parser:"@@"`
	Minimum    Uint `parser:"@@ ')'"`
}

type NSRecord struct {
	Pos        lexer.Position
	NameServer Name `parser:"'@' 'IN' 'NS' @@"`
}

type ARecord struct {
	Pos  lexer.Position
	Name Name `parser:"@@"`
	Ip   Ipv4 `parser:"'IN' 'A' @@"`
}

type MXRecord struct {
	Pos         lexer.Position
	Priority    Uint `parser:"'@' 'IN' 'MX' @@"`
	EmailServer Name `parser:"@@"`
}

type TXTRecord struct {
	Pos   lexer.Position
	Value Text `parser:"'@' 'IN' 'TXT' @@"`
}

type Domain struct {
	Pos   lexer.Position
	Value string `parser:"@Domain"`
}

type Ttl struct {
	Pos   lexer.Position
	Value string `parser:"@Ttl"`
}

type Name struct {
	Pos   lexer.Position
	Value string `parser:"@Name"`
}

type Ipv4 struct {
	Pos   lexer.Position
	Value string `parser:"@Ipv4"`
}

type Uint struct {
	Pos   lexer.Position
	Value uint `parser:"@Uint"`
}

type Text struct {
	Pos   lexer.Position
	Value string `parser:"@String"`
}
