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

var txtRecords = []*models.TXTRecord{
	{Record: models.Record{DomainId: 2}, Value: "value of txt 1"},
	{Record: models.Record{DomainId: 2}, Value: "value of txt 2"},
}

func createTXTRecords() error {
	if err := domains[0].Create(); err != nil {
		return err
	}

	if err := domains[1].Create(); err != nil {
		return err
	}

	for _, txtRecord := range txtRecords {
		if err := txtRecord.Create(); err != nil {
			return err
		}
	}

	return nil
}

func assertEqualTXTRecords(t *testing.T, txt1, txt2 *models.TXTRecord, strict bool) {
	if strict {
		assert.Equal(t, txt1.Id, txt2.Id)
	}
	assert.Equal(t, txt1.DomainId, txt2.DomainId)
	assert.Equal(t, txt1.Value, txt2.Value)
}

func TestListTXTRecords(t *testing.T) {
	tests.WithTestDatabase(
		t, func() {
			if err := createTXTRecords(); err != nil {
				t.Fatal(err)
			}

			router := api.SetupRouter()

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/domains/2/txtRecords", nil)
			router.ServeHTTP(w, req)

			var resp map[string][]*models.TXTRecord
			json.Unmarshal(w.Body.Bytes(), &resp)

			assert.Equal(t, http.StatusOK, w.Code, w.Body.String())
			assert.Len(t, resp["txtRecords"], 2)
			assertEqualTXTRecords(t, txtRecords[1], resp["txtRecords"][1], true)

		},
	)
}

func TestGetTXTRecord(t *testing.T) {
	tests.WithTestDatabase(
		t, func() {
			if err := createTXTRecords(); err != nil {
				t.Fatal(err)
			}

			router := api.SetupRouter()

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/domains/2/txtRecords/2", nil)
			router.ServeHTTP(w, req)

			resp := &models.TXTRecord{}
			json.Unmarshal(w.Body.Bytes(), resp)

			assert.Equal(t, http.StatusOK, w.Code, w.Body.String())
			assertEqualTXTRecords(t, txtRecords[1], resp, true)
		},
	)
}

func TestNewTXTRecord(t *testing.T) {
	tests.WithTestDatabase(
		t, func() {
			if err := createDomains(); err != nil {
				t.Fatal(err)
			}

			router := api.SetupRouter()

			jsonData, _ := json.Marshal(txtRecords[0])

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/domains/2/txtRecords", bytes.NewBuffer(jsonData))
			router.ServeHTTP(w, req)

			resp := &models.TXTRecord{}
			json.Unmarshal(w.Body.Bytes(), &resp)

			assert.Equal(t, http.StatusCreated, w.Code, w.Body.String())
			assertEqualTXTRecords(t, txtRecords[0], resp, false)
		},
	)
}

func TestUpdateTXTRecord(t *testing.T) {
	tests.WithTestDatabase(
		t, func() {
			if err := createTXTRecords(); err != nil {
				return
			}

			router := api.SetupRouter()

			updatedTXTRecord := models.TXTRecord{
				Value: "Updated txt value",
			}

			expectedTXTRecord := &models.TXTRecord{
				Record: models.Record{
					Id:       2,
					DomainId: 2,
				},
				Value: "Updated txt value",
			}

			jsonData, _ := json.Marshal(updatedTXTRecord)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PATCH", "/api/domains/2/txtRecords/2", bytes.NewBuffer(jsonData))
			router.ServeHTTP(w, req)

			resp := &models.TXTRecord{}
			json.Unmarshal(w.Body.Bytes(), &resp)

			assert.Equal(t, http.StatusOK, w.Code, w.Body.String())
			assertEqualTXTRecords(t, expectedTXTRecord, resp, true)
		},
	)
}

func TestDeleteTXTRecord(t *testing.T) {
	tests.WithTestDatabase(
		t, func() {
			if err := createTXTRecords(); err != nil {
				t.Fatal(err)
			}

			router := api.SetupRouter()

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", "/api/domains/2/txtRecords/2", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusNoContent, w.Code, w.Body.String())
			assert.Empty(t, w.Body.String())

			var count int64
			models.DB.Model(
				&models.TXTRecord{},
			).Where(
				&models.TXTRecord{Record: models.Record{Id: 2, DomainId: 2}},
			).Count(&count)
			assert.Equal(t, int64(0), count)
		},
	)
}
