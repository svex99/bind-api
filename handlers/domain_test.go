package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/svex99/bind-api/api"
	"github.com/svex99/bind-api/internal/tests"
	"github.com/svex99/bind-api/models"
)

var domains = []*models.Domain{
	{Name: "domain1.com", NameServer: "ns1", NSIp: "1.1.1.1", Ttl: "1d"},
	{Name: "domain2.com", NameServer: "ns2", NSIp: "2.2.2.2", Ttl: "2d"},
	{Name: "domain3.com", NameServer: "ns3", NSIp: "3.3.3.3", Ttl: "3d"},
}

func createDomains() {
	for _, domain := range domains {
		domain.Create()
		// models.DB.Create(domain)
	}
}

func assertEqualDomains(t *testing.T, d1, d2 *models.Domain) {
	assert.Equal(t, d1.Id, d2.Id)
	assert.Equal(t, d1.Name, d2.Name)
	assert.Equal(t, d1.NameServer, d2.NameServer)
	assert.Equal(t, d1.NSIp, d2.NSIp)
	assert.Equal(t, d1.Ttl, d2.Ttl)
}

func getDomainFromDB(domainId uint) *models.Domain {
	domain := &models.Domain{Id: domainId}

	if err := models.DB.First(domain).Error; err != nil {
		return nil
	}

	return domain
}

func TestListDomains(t *testing.T) {
	tests.WithTestDatabase(
		t, func() error {
			createDomains()

			router := api.SetupRouter()

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/domains", nil)
			router.ServeHTTP(w, req)

			var resp map[string][]*models.Domain
			json.Unmarshal(w.Body.Bytes(), &resp)

			assert.Equal(t, http.StatusOK, w.Code, w.Body.String())
			assert.Len(t, resp["domains"], 3)
			assertEqualDomains(t, domains[0], resp["domains"][0])

			return nil
		},
	)
}

func TestGetDomain(t *testing.T) {
	tests.WithTestDatabase(
		t, func() error {
			createDomains()

			router := api.SetupRouter()

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/domains/1", nil)
			router.ServeHTTP(w, req)

			resp := &models.Domain{}
			json.Unmarshal(w.Body.Bytes(), resp)

			assert.Equal(t, http.StatusOK, w.Code, w.Body.String())
			assertEqualDomains(t, domains[0], resp)

			return nil
		},
	)
}

func TestNewDomain(t *testing.T) {
	tests.WithTestDatabase(
		t, func() error {
			router := api.SetupRouter()

			jsonData, _ := json.Marshal(domains[0])

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/domains", bytes.NewBuffer(jsonData))
			router.ServeHTTP(w, req)

			resp := &models.Domain{}
			json.Unmarshal(w.Body.Bytes(), &resp)

			assert.Equal(t, http.StatusCreated, w.Code, w.Body.String())
			assertEqualDomains(t, domains[0], resp)

			assertEqualDomains(t, domains[0], getDomainFromDB(1))

			return nil
		},
	)
}

func TestUpdateDomain(t *testing.T) {
	tests.WithTestDatabase(
		t, func() error {
			createDomains()

			router := api.SetupRouter()

			updatedDomain := models.UpdateDomainForm{
				Name:       "new-name.com",
				NameServer: "new-name-server",
				NSIp:       "123.123.123.123",
				Ttl:        "10d",
			}

			expectedDomain := &models.Domain{
				Id:         1,
				Name:       updatedDomain.Name,
				NameServer: updatedDomain.NameServer,
				NSIp:       updatedDomain.NSIp,
				Ttl:        updatedDomain.Ttl,
			}

			jsonData, _ := json.Marshal(updatedDomain)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PATCH", "/api/domains/1", bytes.NewBuffer(jsonData))
			router.ServeHTTP(w, req)

			resp := &models.Domain{}
			json.Unmarshal(w.Body.Bytes(), &resp)

			assert.Equal(t, http.StatusOK, w.Code, w.Body.String())
			assertEqualDomains(t, expectedDomain, resp)

			assertEqualDomains(t, expectedDomain, getDomainFromDB(1))

			return nil
		},
	)
}

func TestDeleteDomain(t *testing.T) {
	tests.WithTestDatabase(
		t, func() error {
			createDomains()

			router := api.SetupRouter()

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", "/api/domains/1", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusNoContent, w.Code, w.Body.String())
			assert.Empty(t, w.Body.String())

			var count int64
			models.DB.Model(&models.Domain{}).Where(&models.Domain{Id: 1}).Count(&count)
			assert.Equal(t, int64(0), count)

			return nil
		},
	)
}
