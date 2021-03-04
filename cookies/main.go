package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// https://github.com/gorilla/securecookie
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/cookie", getCookie).Methods(http.MethodGet)
	router.HandleFunc("/restricted", getRestricted).Methods(http.MethodGet)

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

func getCookie(w http.ResponseWriter, r *http.Request) {
	// todo: sign in
	c := http.Cookie{
		Name:       "session_id",
		Value:      "1111111111",
		Domain:     "localhost",
		Secure:     true,
		HttpOnly:   true,
		SameSite:   http.SameSiteStrictMode,
	}
	http.SetCookie(w, &c)
}

func getRestricted(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sid := c.Value
	if sid != "1111111111" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Write([]byte("ok"))
}