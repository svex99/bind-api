package setting

import (
	"log"

	"github.com/go-ini/ini"
)

type AppSetting struct {
	JwtSecret         string
	TokenHourLifespan int
	BindCfgPath       string
	BindDbsPath       string
	BindAdmin         string
}

var App = &AppSetting{}

var cfg *ini.File

func init() {
	var err error
	cfg, err = ini.Load("app.ini")
	if err != nil {
		log.Fatalf("Error parsing app.ini file: %v", err)
	}

	mapTo("app", App)
}

func mapTo(section string, v interface{}) {
	if err := cfg.Section(section).MapTo(v); err != nil {
		log.Fatalf("Error mapping section %s in .ini file: %v", section, err)
	}
}
