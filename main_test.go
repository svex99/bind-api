package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime/pprof"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/svex99/bind-api/api"
	"github.com/svex99/bind-api/internal/tests"
	"github.com/svex99/bind-api/schemas"
	"github.com/svex99/bind-api/services/bind"
	"github.com/svex99/bind-api/services/bind/parser"
)

func BenchmarkAddZoneLatency(b *testing.B) {
	amounts := []int{0, 250, 500, 1000, 2000, 4000}

	for _, amount := range amounts {
		requestsData := [][]byte{}
		for i := 1; i <= 100; i++ {
			body, _ := json.Marshal(schemas.ZoneData{
				Origin:     fmt.Sprintf("test-api%d.com", i),
				Ttl:        "2d",
				NameServer: "test-ns",
				Admin:      "svex",
				Refresh:    1234,
				Retry:      5678,
				Expire:     9012,
				Minimum:    3456,
			})
			requestsData = append(requestsData, body)
		}

		b.Run(fmt.Sprintf("AddZonesLatencyWith%d", amount), func(b *testing.B) {
			for _, body := range requestsData {
				b.StopTimer()
				if err := tests.CreateZonesBulk(amount); err != nil {
					b.Fatal(err)
				}

				router := api.SetupRouter(false)
				bind.Service.Init()

				w := httptest.NewRecorder()

				req, err := http.NewRequest("POST", "/api/zones", bytes.NewReader(body))
				if err != nil {
					b.Fatal(err)
				}

				b.StartTimer()
				router.ServeHTTP(w, req)

				b.StopTimer()
				if w.Code != http.StatusCreated {
					b.Fatalf(
						"invalid status code: got %d, expected %d\nbody:%v",
						w.Code, http.StatusCreated, w.Body,
					)
				}
				b.StartTimer()
			}
		})
	}
}

func BenchmarkGetZonesLatency(b *testing.B) {
	amounts := []int{250, 500, 1000, 2000, 4000}

	for _, amount := range amounts {
		if err := tests.CreateZonesBulk(amount); err != nil {
			b.Fatal(err)
		}

		router := api.SetupRouter(false)
		bind.Service.Init()

		b.Run(fmt.Sprintf("GetZonesLatencyWith%d", amount), func(b *testing.B) {
			w := httptest.NewRecorder()

			req, err := http.NewRequest("GET", "/api/zones", nil)
			if err != nil {
				b.Fatal(err)
			}

			b.ResetTimer()
			router.ServeHTTP(w, req)
			b.StopTimer()

			if w.Code != http.StatusOK {
				b.Fatalf(
					"invalid status code: got %d, expected %d\nbody:%v",
					w.Code, http.StatusCreated, w.Body,
				)
			}
		})
	}
}

func BenchmarkBindServiceMemory(b *testing.B) {
	amounts := []int{0, 250, 500, 1000, 2000, 4000, 8000}

	for _, amount := range amounts {
		tests.CreateZonesBulk(amount)

		b.Run(fmt.Sprintf("BenchmarkBindServiceMemoryWith%d", amount), func(b *testing.B) {
			memProfile, _ := os.Create(fmt.Sprintf("mem_%d.out", amount))

			bind.Service.Init()

			pprof.WriteHeapProfile(memProfile)
		})
	}
}

func TestAPIFlow(t *testing.T) {
	testData := []struct {
		origin string
		ns     string
		ip     string
		nsd    string
	}{
		{"flow-test-1.com", "ns1", "10.1.3.1", "ns1.flow-test-1.com"},
		{"flow-test-2.com", "ns2", "10.1.3.2", "ns2.flow-test-2.com"},
		{"flow-test-3.com", "ns3", "10.1.3.3", "ns3.flow-test-3.com"},
		{"flow-test-4.com", "ns4", "10.1.3.4", "ns4.flow-test-4.com"},
		{"flow-test-5.com", "ns5", "10.1.3.5", "ns5.flow-test-5.com"},
	}

	if err := tests.CleanBindConfig(); err != nil {
		t.Fatal(err)
	}

	router := api.SetupRouter(false)
	bind.Service.Init()

	t.Run("TestAddZone", func(t *testing.T) {
		for _, td := range testData {
			// Add zone
			body, _ := json.Marshal(schemas.ZoneData{
				Origin:     td.origin,
				Ttl:        "2d",
				NameServer: td.ns,
				Admin:      "svex",
				Refresh:    86400,
				Retry:      7200,
				Expire:     3600000,
				Minimum:    172800,
			})
			w := httptest.NewRecorder()
			req, err := http.NewRequest("POST", "/api/zones", bytes.NewReader(body))
			if err != nil {
				t.Fatal(err)
			}
			router.ServeHTTP(w, req)
			if w.Code != http.StatusCreated {
				t.Fatalf(
					"invalid status code: got %d, expected %d\nbody:%v",
					w.Code, http.StatusCreated, w.Body,
				)
			}

			// Add NS record
			tests.Serve(
				router, "POST", "/api/zones/"+td.origin+"/records",
				parser.NSRecord{Type: "NS", NameServer: td.ns},
				http.StatusCreated,
			)

			// Add A record for NS
			tests.Serve(
				router, "POST", "/api/zones/"+td.origin+"/records",
				parser.ARecord{Name: td.ns, Type: "A", Ip: td.ip},
				http.StatusCreated,
			)

			ips, err := net.LookupIP(td.nsd)
			if err != nil {
				t.Fatal(err)
			}

			assert.Len(t, ips, 1)
			assert.Equal(t, td.ip, ips[0].String())
		}
	})

	t.Run("TestMXRecord", func(t *testing.T) {
		data := testData[0]
		endpoint := "/api/zones/" + data.origin + "/records"
		mxData := []struct {
			Type        string `json:"type"`
			Priority    string `json:"priority"`
			EmailServer string `json:"emailServer"`
		}{{"MX", "101", "email-server1"}, {"MX", "102", "email-server2"}, {"MX", "103", "email-server3"}}

		t.Run("TestAddMX", func(t *testing.T) {
			for _, mxD := range mxData {
				if err := tests.Serve(router, "POST", endpoint, mxD, http.StatusCreated); err != nil {
					t.Fatal(err)
				}
			}

			mxRecords, err := net.LookupMX(data.origin)
			if err != nil {
				t.Fatal(err)
			}

			assert.Len(t, mxRecords, 3)
			assert.Equal(t, "email-server1."+data.origin+".", mxRecords[0].Host)
			assert.Equal(t, uint16(101), mxRecords[0].Pref)
			assert.Equal(t, "email-server3."+data.origin+".", mxRecords[2].Host)
			assert.Equal(t, uint16(103), mxRecords[2].Pref)
		})

		t.Run("TestDelMX", func(t *testing.T) {
			if err := tests.Serve(router, "DELETE", endpoint, mxData[1], http.StatusOK); err != nil {
				t.Fatal(err)
			}

			mxRecords, err := net.LookupMX(data.origin)
			if err != nil {
				t.Fatal(err)
			}

			assert.Len(t, mxRecords, 2)
			assert.Equal(t, "email-server1."+data.origin+".", mxRecords[0].Host)
			assert.Equal(t, uint16(101), mxRecords[0].Pref)
			assert.Equal(t, "email-server3."+data.origin+".", mxRecords[1].Host)
			assert.Equal(t, uint16(103), mxRecords[1].Pref)

		})
	})

	t.Run("TestCNAMERecord", func(t *testing.T) {
		data := testData[0]
		endpoint := "/api/zones/" + data.origin + "/records"

		t.Run("TestAddCNAME", func(t *testing.T) {
			// Add destination A record
			if err := tests.Serve(
				router, "POST", endpoint,
				parser.ARecord{Name: "cname-dest", Type: "A", Ip: "7.7.7.7"},
				http.StatusCreated,
			); err != nil {
				t.Fatal(err)
			}

			// Add CNAME record
			if err := tests.Serve(
				router, "POST", endpoint,
				parser.CNAMERecord{SrcName: "cname-source", Type: "CNAME", DstName: "cname-dest"},
				http.StatusCreated,
			); err != nil {
				t.Fatal(err)
			}

			cname, err := net.LookupCNAME("cname-source." + data.origin)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, "cname-dest."+data.origin+".", cname)
		})

		t.Run("TestDelCNAME", func(t *testing.T) {
			if err := tests.Serve(
				router, "DELETE", endpoint,
				parser.CNAMERecord{SrcName: "cname-source", Type: "CNAME", DstName: "cname-dest"},
				http.StatusOK,
			); err != nil {
				t.Fatal(err)
			}

			_, err := net.LookupCNAME("cname-source." + data.origin)

			assert.NotNil(t, err)
		})
	})

	t.Run("TestTXTRecord", func(t *testing.T) {
		data := testData[0]
		endpoint := "/api/zones/" + data.origin + "/records"

		txtData := []string{"key1=value-body", "key2=value-body", "key3=value-body"}

		t.Run("TestAddTXT", func(t *testing.T) {
			for _, txtD := range txtData {
				if err := tests.Serve(
					router, "POST", endpoint,
					parser.TXTRecord{Type: "TXT", Value: txtD},
					http.StatusCreated,
				); err != nil {
					t.Fatal(err)
				}
			}

			txts, err := net.LookupTXT(data.origin)
			if err != nil {
				t.Fatal(err)
			}

			assert.Len(t, txts, 3)
			assert.ElementsMatch(t, txtData, txts)
		})

		t.Run("TestDelTXT", func(t *testing.T) {
			for _, txtD := range txtData[1:] {
				if err := tests.Serve(
					router, "DELETE", endpoint,
					parser.TXTRecord{Type: "TXT", Value: txtD},
					http.StatusOK,
				); err != nil {
					t.Fatal(err)
				}
			}

			txts, err := net.LookupTXT(data.origin)
			if err != nil {
				t.Fatal(err)
			}

			assert.Len(t, txts, 1)
			assert.Equal(t, txtData[0], txts[0])
		})
	})

	t.Run("TestDelZone", func(t *testing.T) {
		for _, td := range testData[1:] {
			if err := tests.Serve(
				router, "DELETE", "/api/zones/"+td.origin, "none", http.StatusNoContent,
			); err != nil {
				t.Fatal(err)
			}

			_, err := net.LookupIP(td.nsd)

			assert.NotNil(t, err)
		}
	})
}
