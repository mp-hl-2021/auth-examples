package main

import (
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go/request"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

// openssl genrsa -out app.rsa keysize
// openssl rsa -in app.rsa -pubout > app.rsa.pub

var (
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
)

func init() {
	privateBytes, err := ioutil.ReadFile("app.rsa")
	if err != nil {
		panic(err)
	}
	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateBytes)
	if err != nil {
		panic(err)
	}
	publicBytes, err := ioutil.ReadFile("app.rsa.pub")
	if err != nil {
		panic(err)
	}
	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicBytes)
	if err != nil {
		panic(err)
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/signin", postSignin).Methods(http.MethodPost)
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

func postSignin(w http.ResponseWriter, r *http.Request) {
	ts, err := createToken("bob")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/jwt")
	w.WriteHeader(http.StatusOK)

	fmt.Fprintln(w, ts) // todo: what is the matter? Why w.Write failed me?
	w.Write([]byte(ts))
}

// curl -v -H 'Accept: application/json' -H "Authorization: Bearer ${TOKEN}" localhost:8080/restricted
func getRestricted(w http.ResponseWriter, r *http.Request) {
	t, err := request.ParseFromRequest(r, request.OAuth2Extractor, func(t *jwt.Token) (interface{}, error) {
		return publicKey, nil
	}, request.WithClaims(&CustomClaims{}))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		panic(err)
		return
	}
	w.Write([]byte("hi, " + t.Claims.(*CustomClaims).Username))
	return
}

type AccountInfo struct {
	Username string
}

type CustomClaims struct {
	*jwt.StandardClaims
	AccountInfo
}

func createToken(username string) (string, error) {
	t := jwt.New(jwt.GetSigningMethod("RS256"))
	t.Claims = &CustomClaims{
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 10).Unix(),
		},
		AccountInfo: AccountInfo{username},
	}
	return t.SignedString(privateKey)
}
