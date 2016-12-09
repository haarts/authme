package main

import (
	"fmt"
	"log"
	"net/http"
)

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // TODO this returns 'plain/text' which is wrong
		return
	}
	if r.PostFormValue("username") == "" {
		http.Error(w, "'username' missing", http.StatusBadRequest)
		return
	}
	if r.PostFormValue("password") == "" {
		http.Error(w, "'password' missing", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello")
}

func authenticatedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello")
}

func resetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello")
}

func main() {
	http.Handle("/register", http.HandlerFunc(registerHandler))
	http.Handle("/login", http.HandlerFunc(loginHandler))
	http.Handle("/authenticated", http.HandlerFunc(authenticatedHandler))
	http.Handle("/reset", http.HandlerFunc(resetHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
