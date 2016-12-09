package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterHandler(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodPost,
		"http://example.com/foo",
		strings.NewReader(url.Values{"username": {"bar"}, "password": {"bla"}}.Encode()),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	w := httptest.NewRecorder()
	registerHandler(w, req)

	assert.Equal(t, 201, w.Code)
}

func TestRegisterHandlerWithoutUsername(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodPost,
		"http://example.com/foo",
		nil,
	)
	w := httptest.NewRecorder()
	registerHandler(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestRegisterHandlerWithoutPassword(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodPost,
		"http://example.com/foo",
		strings.NewReader(url.Values{"username": {"bar"}}.Encode()),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	w := httptest.NewRecorder()
	registerHandler(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "'password' missing\n", w.Body.String())
}
