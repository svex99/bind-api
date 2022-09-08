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

var emails = []*models.Email{
	{Priority: 100, Name: "email1", Ip: "100.100.100.100"},
	{Priority: 200, Name: "email2", Ip: "200.200.200.200"},
}

func createEmails() error {
	if err := domains[0].Create(); err != nil {
		return err
	}

	if err := domains[1].Create(); err != nil {
		return err
	}

	for _, email := range emails {
		if err := email.Create(2); err != nil {
			return err
		}
	}

	return nil
}

func assertEqualEmails(t *testing.T, e1, e2 *models.Email, strict bool) {
	if strict {
		assert.Equal(t, e1.Id, e2.Id)
	}
	assert.Equal(t, e1.Priority, e2.Priority)
	assert.Equal(t, e1.Name, e2.Name)
	assert.Equal(t, e1.Ip, e2.Ip)
}

func TestListEmails(t *testing.T) {
	tests.SetupTestDatabase(t)

	if err := createEmails(); err != nil {
		t.Fatal(err)
	}

	router := api.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/domains/2/emails", nil)
	router.ServeHTTP(w, req)

	var resp map[string][]*models.Email
	json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Equal(t, http.StatusOK, w.Code, w.Body.String())
	assert.Len(t, resp["emails"], 2)
	assertEqualEmails(t, emails[1], resp["emails"][1], true)
}

func TestGetEmail(t *testing.T) {
	tests.SetupTestDatabase(t)

	if err := createEmails(); err != nil {
		t.Fatal(err)
	}

	router := api.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/domains/2/emails/2", nil)
	router.ServeHTTP(w, req)

	resp := &models.Email{}
	json.Unmarshal(w.Body.Bytes(), resp)

	assert.Equal(t, http.StatusOK, w.Code, w.Body.String())
	assertEqualEmails(t, emails[1], resp, true)
}

func TestNewEmail(t *testing.T) {
	tests.SetupTestDatabase(t)

	if err := createDomains(); err != nil {
		t.Fatal(err)
	}

	router := api.SetupRouter()

	jsonData, _ := json.Marshal(emails[0])

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/domains/2/emails", bytes.NewBuffer(jsonData))
	router.ServeHTTP(w, req)

	resp := &models.Email{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Equal(t, http.StatusCreated, w.Code, w.Body.String())
	assertEqualEmails(t, emails[0], resp, false)
}

func TestUpdateEmail(t *testing.T) {
	tests.SetupTestDatabase(t)

	if err := createEmails(); err != nil {
		return
	}

	router := api.SetupRouter()

	updatedEmail := models.UpdateEmailForm{
		Priority: 500,
		Name:     "new-subdomain-name",
		Ip:       "123.123.123.123",
	}

	expectedEmail := &models.Email{
		Id:       2,
		Priority: updatedEmail.Priority,
		Name:     updatedEmail.Name,
		Ip:       updatedEmail.Ip,
	}

	jsonData, _ := json.Marshal(updatedEmail)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/api/domains/2/emails/2", bytes.NewBuffer(jsonData))
	router.ServeHTTP(w, req)

	resp := &models.Email{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Equal(t, http.StatusOK, w.Code, w.Body.String())
	assertEqualEmails(t, expectedEmail, resp, true)
}

func TestDeleteEmail(t *testing.T) {
	tests.SetupTestDatabase(t)

	if err := createEmails(); err != nil {
		t.Fatal(err)
	}

	router := api.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/domains/2/emails/2", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code, w.Body.String())
	assert.Empty(t, w.Body.String())

	var count int64
	models.DB.Model(
		&models.MXRecord{},
	).Where(
		&models.MXRecord{Record: models.Record{Id: 2, DomainId: 2}},
	).Count(&count)
	assert.Equal(t, int64(0), count)
}
