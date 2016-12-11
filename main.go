package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pzduniak/argon2"
)

// FIXME why export this?
type App struct {
	db *sql.DB
}

func secureSalt() ([]byte, error) {
	salt := make([]byte, 64)

	_, err := rand.Read(salt)
	if err != nil {
		return salt, err
	}

	return salt, nil
}

func encryptPassword(salt []byte, password string) ([]byte, error) {
	// TODO check out these parameters
	encryptedPassword, err := argon2.Key([]byte(password), salt, 13, 4, 4096, 32, argon2.Argon2i)
	if err != nil {
		return []byte{}, err
	}

	return encryptedPassword, nil
}

func (a *App) storeUser(username, password string) error {
	salt, err := secureSalt()
	if err != nil {
		return err
	}

	encryptedPassword, err := encryptPassword(salt, password)
	if err != nil {
		return err
	}

	_, err = a.db.Exec(
		"INSERT INTO users (username, encrypted_password, salt) VALUES ($1, $2, $3)",
		username,
		hex.EncodeToString(encryptedPassword),
		hex.EncodeToString(salt),
	)

	return err
}

func usernameAndPasswordFromForm(r *http.Request) (string, string, error) {
	if err := r.ParseForm(); err != nil {
		return "", "", err
	}
	if r.PostFormValue("username") == "" {
		return "", "", errors.New("'username' must be present")
	}
	if r.PostFormValue("password") == "" {
		return "", "", errors.New("'password' must be present")
	}

	return r.PostFormValue("username"), r.PostFormValue("password"), nil

}

func (a *App) registerHandler(w http.ResponseWriter, r *http.Request) {
	username, password, err := usernameAndPasswordFromForm(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.storeUser(username, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *App) storeSession() (string, error) {
	session := make([]byte, 16*8)
	_, err := rand.Read(session)
	if err != nil {
		return "", err
	}

	hexedSession := hex.EncodeToString(session)
	_, err = a.db.Exec("INSERT INTO sessions(session) VALUES(?)", hexedSession)
	if err != nil {
		return "", err
	}

	return hexedSession, nil
}

func (a *App) loginHandler(w http.ResponseWriter, r *http.Request) {
	username, password, err := usernameAndPasswordFromForm(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var salt, storedEncryptedPassword string
	err = a.db.QueryRow("SELECT salt, encrypted_password FROM users WHERE username=?", username).Scan(&salt, &storedEncryptedPassword)
	switch {
	case err == sql.ErrNoRows:
		http.Error(w, "login failed", http.StatusUnauthorized)
		return
	case err != nil:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	decodedSalt, err := hex.DecodeString(salt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	generatedEncryptedPassword, err := encryptPassword(decodedSalt, password)
	if hex.EncodeToString(generatedEncryptedPassword) != storedEncryptedPassword {
		http.Error(w, "login failed", http.StatusUnauthorized)
		return
	}

	session, err := a.storeSession()

	http.SetCookie(w, &http.Cookie{
		Name:     "sessionid",
		Value:    session,
		Expires:  time.Now().AddDate(0, 0, 7),
		Secure:   true,
		HttpOnly: true,
		// TODO set domain
	})
}

func (a *App) authenticatedHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("sessionid")
	if err != nil {
		http.Error(w, "invalid session", http.StatusUnauthorized)
		return
	}

	var session string
	err = a.db.QueryRow("SELECT session FROM sessions WHERE session = ?", cookie.Value).Scan(&session)
	switch {
	case err == sql.ErrNoRows:
		http.Error(w, "invalid session", http.StatusUnauthorized)
		return
	case err != nil:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func resetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello")
}

func main() {
	db, err := sql.Open("sqlite3", "./users.db")
	if err != nil {
		fmt.Printf("err = %+v\n", err)
		return
	}
	defer db.Close()

	app := App{
		db: db,
	}

	http.Handle("/register", http.HandlerFunc(app.registerHandler))
	http.Handle("/login", http.HandlerFunc(app.loginHandler))
	http.Handle("/authenticated", http.HandlerFunc(app.authenticatedHandler))
	http.Handle("/reset", http.HandlerFunc(resetHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
