package parser_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/svex99/bind-api/services/bind/parser"
)

func TestParser(t *testing.T) {
	file, err := os.Open("testdata/named.conf.local")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(parser.ConfParser.String())

	conf, err := parser.ConfParser.Parse("", file)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "example.com", conf.Zones[0].Name)
	assert.Equal(t, "master", conf.Zones[0].Type)
	assert.Equal(t, "/var/lib/bind/db.example.com", conf.Zones[0].File)

	assert.Len(t, conf.Zones, 2)

	content, err := os.ReadFile("testdata/named.conf.local")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(conf.String())
	assert.Equal(t, conf.String(), string(content))
}
