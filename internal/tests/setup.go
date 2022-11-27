package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/svex99/bind-api/pkg/setting"
	"github.com/svex99/bind-api/services/bind/parser"
)

func CleanBindConfig() error {
	dir, err := os.ReadDir(setting.Bind.LibPath)
	if err != nil {
		return err
	}

	for _, d := range dir {
		if err := os.RemoveAll(path.Join([]string{setting.Bind.LibPath, d.Name()}...)); err != nil {
			return err
		}
	}

	if err := os.Truncate(setting.Bind.ConfPath+"named.conf.local", 0); err != nil {
		return err
	}

	return nil
}

func CreateZonesBulk(amount int) error {
	if err := CleanBindConfig(); err != nil {
		return err
	}

	bindConf := parser.BindConf{
		Zones: []*parser.Zone{},
	}

	for i := 1; i <= amount; i++ {
		zoneName := fmt.Sprintf("bulk-domain-%d.com", i)
		ns := fmt.Sprintf("ns-bulk-domain-%d", i)

		bindConf.Zones = append(bindConf.Zones, &parser.Zone{
			Name: zoneName,
			Type: "primary",
			File: "/var/lib/bind/db." + zoneName,
		})

		dc := parser.ZoneConf{
			Origin: zoneName,
			Ttl:    "20d",
			SOARecord: &parser.SOARecord{
				NameServer: ns,
				Admin:      "svex",
				Serial:     11111,
				Refresh:    22222,
				Retry:      33333,
				Expire:     44444,
				Minimum:    55555,
			},
			Records: []parser.Record{
				parser.NSRecord{Type: "NS", NameServer: ns},
				parser.ARecord{Name: ns, Type: "A", Ip: "123.123.123.123"},
			},
		}

		if _, err := dc.WriteToDisk(setting.Bind.LibPath + "db." + zoneName); err != nil {
			return err
		}
	}

	if _, err := bindConf.WriteToDisk(setting.Bind.ConfPath + "named.conf.local"); err != nil {
		return err
	}

	return nil
}

func Serve(router *gin.Engine, method string, url string, data any, expectedCode int) error {
	var body *bytes.Reader = nil

	if data != nil {
		encodedBody, err := json.Marshal(data)
		if err != nil {
			return err
		}
		body = bytes.NewReader(encodedBody)
	}

	w := httptest.NewRecorder()

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}

	router.ServeHTTP(w, req)

	if w.Code != expectedCode {
		return fmt.Errorf(
			"invalid status code: got %d, expected %d\nbody:%v",
			w.Code, http.StatusCreated, w.Body,
		)
	}

	return nil
}
