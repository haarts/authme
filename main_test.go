package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoginHandler(t *testing.T) {
	db, err := initializeDatabase()
	require.NoError(t, err)

	app := App{db: db}
	app.storeUser("foo", "bar")

	req := httptest.NewRequest(
		http.MethodPost,
		"http://example.com/foo",
		strings.NewReader(url.Values{"username": {"foo"}, "password": {"bar"}}.Encode()),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	w := httptest.NewRecorder()
	app.loginHandler(w, req)

	assert.Equal(t, 200, w.Code)

	var count int
	require.NoError(t, db.QueryRow("SELECT COUNT(*) FROM sessions").Scan(&count))
	assert.Equal(t, 1, count)
}

func TestLoginFailed(t *testing.T) {
	db, err := initializeDatabase()
	require.NoError(t, err)

	app := App{db: db}
	app.storeUser("foo", "bar")

	req := httptest.NewRequest(
		http.MethodPost,
		"http://example.com/foo",
		strings.NewReader(url.Values{"username": {"TOTALLY DIFFERENT"}, "password": {"bar"}}.Encode()),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	w := httptest.NewRecorder()
	app.loginHandler(w, req)

	fmt.Printf("w = %+v\n", w)
	assert.Equal(t, 401, w.Code)
}

func TestRegisterHandler(t *testing.T) {
	db, err := initializeDatabase()
	require.NoError(t, err)

	app := App{db: db}

	req := httptest.NewRequest(
		http.MethodPost,
		"http://example.com/foo",
		strings.NewReader(url.Values{"username": {"bar"}, "password": {"bla"}}.Encode()),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	w := httptest.NewRecorder()
	app.registerHandler(w, req)

	assert.Equal(t, 201, w.Code)
	assert.Equal(t, "", w.Body.String())
	var count int
	require.NoError(t, db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count))
	assert.Equal(t, 1, count)
}

func TestRegisterHandlerWithoutUsername(t *testing.T) {
	app := App{}

	req := httptest.NewRequest(
		http.MethodPost,
		"http://example.com/foo",
		nil,
	)
	w := httptest.NewRecorder()
	app.registerHandler(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestRegisterHandlerWithoutPassword(t *testing.T) {
	app := App{}

	req := httptest.NewRequest(
		http.MethodPost,
		"http://example.com/foo",
		strings.NewReader(url.Values{"username": {"bar"}}.Encode()),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	w := httptest.NewRecorder()
	app.registerHandler(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "'password' must be present\n", w.Body.String())
}

func initializeDatabase() (*sql.DB, error) {
	schema, err := ioutil.ReadFile("schema.sql")
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(string(schema))
	if err != nil {
		return nil, err
	}

	return db, nil
}
