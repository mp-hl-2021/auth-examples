package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/ok", authorize(getOK)).Methods(http.MethodGet)
	router.HandleFunc("/test", authorize(getTest)).Methods(http.MethodGet)

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

func getOK(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK!"))
}

func getTest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Test"))
}

func authorize(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if !ok { // check authorization header
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if u != "test" || p != "test" { // check user credentials
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		handler(w, r)
	}
}