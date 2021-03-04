package main

import (
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gorilla/mux"
)

type Account struct {
	Username string
	Password []byte
}

var Accounts map[string]Account

func main() {
	Accounts = make(map[string]Account)

	router := mux.NewRouter()
	router.HandleFunc("/signup", postSignup).Methods(http.MethodPost)
	router.HandleFunc("/signin", postSignin).Methods(http.MethodPost)

	server := http.Server{
		Addr:         "localhost:8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,

		Handler: router,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func postSignup(w http.ResponseWriter, r *http.Request) {
	c := Credentials{}
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// todo: validate credentials!
	if _, ok := Accounts[c.Username]; ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// todo: remember to add throttling.
	h, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	Accounts[c.Username] = Account{
		Username: c.Username,
		Password: h,
	}
}

func postSignin(w http.ResponseWriter, r *http.Request) {
	c := Credentials{}
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// todo: validate
	a, ok := Accounts[c.Username]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := bcrypt.CompareHashAndPassword(a.Password, []byte(c.Password)); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write([]byte("Successfully logged in"))
}