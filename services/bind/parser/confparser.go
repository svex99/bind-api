package parser

import (
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/svex99/bind-api/pkg/file"
)

var (
	ConfLexer = lexer.MustSimple([]lexer.SimpleRule{
		{Name: "Keyword", Pattern: `[a-zA-Z\-]+`},
		{Name: "String", Pattern: `"[^"\n]*"`},
		{Name: "Punct", Pattern: `[\{\}\;]`},
		{Name: "Comment", Pattern: `//[^\n]*\n+`},
		{Name: "Whitespace", Pattern: `[ \t\r\n]+`},
	})
	ConfParser = participle.MustBuild[BindConf](
		participle.Lexer(ConfLexer),
		participle.Unquote("String"),
		participle.Elide("Whitespace", "Comment"),
		participle.UseLookahead(1),
	)
)

type BindConf struct {
	Zones []*Zone `parser:"@@*"`
}

type Zone struct {
	Name string `parser:"'zone' @String '{'"`
	Type string `parser:"'type' @('primary'|'master') ';'"`
	File string `parser:"'file' @String ';' '}' ';'"`
}

func (bc *BindConf) WriteToDisk(filename string) (func(), error) {
	// Create a backup of config if file exists
	rollback := file.MakeBackup(filename)

	if err := os.WriteFile(filename, []byte(bc.String()), 0666); err != nil {
		return rollback, err
	}

	return rollback, nil
}

func (bc *BindConf) AddZone(dc *ZoneConf) error {
	for _, zone := range bc.Zones {
		if zone.Name == dc.Origin {
			return fmt.Errorf("zone already exists")
		}
	}

	bc.Zones = append(bc.Zones, &Zone{dc.Origin, "master", "/var/lib/bind/db." + dc.Origin})

	return nil
}

func (bc *BindConf) DeleteZone(dc *ZoneConf) error {
	foundIndex := -1

	for i, zone := range bc.Zones {
		if zone.Name == dc.Origin {
			foundIndex = i
			break
		}
	}

	if foundIndex == -1 {
		return fmt.Errorf("zone does not exist")
	}

	bc.Zones = append(bc.Zones[:foundIndex], bc.Zones[foundIndex+1:]...)

	return nil
}

func (bc *BindConf) String() string {
	zones := []string{}
	for _, zone := range bc.Zones {
		zones = append(zones, zone.String())
	}
	return strings.Join(zones, "\n")
}

func (z *Zone) String() string {
	return fmt.Sprintf(
		"zone \"%s\" {\n\ttype %s;\n\tfile \"%s\";\n};\n",
		z.Name, z.Type, z.File,
	)
}
