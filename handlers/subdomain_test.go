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

var subdomains = []*models.Subdomain{
	{Name: "subdomain1", Ip: "10.10.10.10"},
	{Name: "subdomain2", Ip: "20.20.20.20"},
}

func createSubdomains() error {
	if err := domains[0].Create(); err != nil {
		return err
	}

	if err := domains[1].Create(); err != nil {
		return err
	}

	for _, subdomain := range subdomains {
		if err := subdomain.Create(2); err != nil {
			return err
		}
	}

	return nil
}

func assertEqualSubdomains(t *testing.T, sd1, sd2 *models.Subdomain, strict bool) {
	if strict {
		assert.Equal(t, sd1.Id, sd2.Id)
	}
	assert.Equal(t, sd1.Name, sd2.Name)
	assert.Equal(t, sd1.Ip, sd2.Ip)
}

func TestListSubdomains(t *testing.T) {
	tests.WithTestDatabase(
		t, func() {
			if err := createSubdomains(); err != nil {
				t.Fatal(err)
			}

			router := api.SetupRouter()

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/domains/2/subdomains", nil)
			router.ServeHTTP(w, req)

			var resp map[string][]*models.Subdomain
			json.Unmarshal(w.Body.Bytes(), &resp)

			assert.Equal(t, http.StatusOK, w.Code, w.Body.String())
			assert.Len(t, resp["subdomains"], 3) // two subdomains plus the name server
			assertEqualSubdomains(t, subdomains[0], resp["subdomains"][1], true)

		},
	)
}

func TestGetSubdomain(t *testing.T) {
	tests.WithTestDatabase(
		t, func() {
			if err := createSubdomains(); err != nil {
				t.Fatal(err)
			}

			router := api.SetupRouter()

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/domains/2/subdomains/4", nil)
			router.ServeHTTP(w, req)

			resp := &models.Subdomain{}
			json.Unmarshal(w.Body.Bytes(), resp)

			assert.Equal(t, http.StatusOK, w.Code, w.Body.String())
			assertEqualSubdomains(t, subdomains[1], resp, true)
		},
	)
}

func TestNewSubdomain(t *testing.T) {
	tests.WithTestDatabase(
		t, func() {
			if err := createDomains(); err != nil {
				t.Fatal(err)
			}
			router := api.SetupRouter()

			jsonData, _ := json.Marshal(subdomains[0])

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/domains/2/subdomains", bytes.NewBuffer(jsonData))
			router.ServeHTTP(w, req)

			resp := &models.Subdomain{}
			json.Unmarshal(w.Body.Bytes(), &resp)

			assert.Equal(t, http.StatusCreated, w.Code, w.Body.String())
			assertEqualSubdomains(t, subdomains[0], resp, false)
		},
	)
}

func TestUpdateSubdomain(t *testing.T) {
	tests.WithTestDatabase(
		t, func() {
			if err := createSubdomains(); err != nil {
				return
			}

			router := api.SetupRouter()

			updatedSubdomain := models.UpdateSubdomainForm{
				Name: "new-subdomain-name",
				Ip:   "123.123.123.123",
			}

			expectedSubdomain := &models.Subdomain{
				Id:   4,
				Name: updatedSubdomain.Name,
				Ip:   updatedSubdomain.Ip,
			}

			jsonData, _ := json.Marshal(updatedSubdomain)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PATCH", "/api/domains/2/subdomains/4", bytes.NewBuffer(jsonData))
			router.ServeHTTP(w, req)

			resp := &models.Subdomain{}
			json.Unmarshal(w.Body.Bytes(), &resp)

			assert.Equal(t, http.StatusOK, w.Code, w.Body.String())
			assertEqualSubdomains(t, expectedSubdomain, resp, true)
		},
	)
}

func TestDeleteSubdomain(t *testing.T) {
	tests.WithTestDatabase(
		t, func() {
			createSubdomains()

			router := api.SetupRouter()

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", "/api/domains/2/subdomains/4", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusNoContent, w.Code, w.Body.String())
			assert.Empty(t, w.Body.String())

			var count int64
			models.DB.Model(
				&models.Domain{},
			).Where(
				&models.ARecord{Record: models.Record{Id: 4, DomainId: 2}},
			).Count(&count)
			assert.Equal(t, int64(0), count)
		},
	)
}
