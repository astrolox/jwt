package jwt_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/pascaldekloe/jwt"
)

// JWTSecret is the HMAC key.
var JWTSecret = []byte("guest")

func ExampleHMACCheckHeader() {
	// request with JWT Bearer
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJsYWthbmUiLCJmbiI6IkxhbmEgQW50aG9ueSBLYW5lIn0.c9qLforxweefXAm7Q8W13N7E6XKT0tNd8rzoablDkx0")

	// execute and print
	rec := httptest.NewRecorder()
	securedHandler(rec, req)
	fmt.Println("HTTP", rec.Result().Status)
	io.Copy(os.Stdout, rec.Result().Body)

	// output: HTTP 200 OK
	// Hello Lana Anthony Kane!
	// You are authorized as lakane.
}

func securedHandler(w http.ResponseWriter, r *http.Request) {
	// verify claims
	claims, err := jwt.HMACCheckHeader(r, JWTSecret)
	if err != nil {
		w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Bearer error="invalid_token", error_description=%q`, err.Error()))
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// verify time constrains
	if !claims.Valid(time.Now()) {
		w.Header().Set("WWW-Authenticate", `Bearer error="invalid_token", error_description="jwt: time contraints exceeded"`)
		http.Error(w, "jwt: time contraints exceeded", http.StatusUnauthorized)
		return
	}

	// verify custom token
	name, ok := claims.String("fn")
	if !ok {
		http.Error(w, `jwt: no "fn" string claim`, http.StatusForbidden)
		return
	}

	fmt.Fprintf(w, "Hello %s!\nYou are authorized as %s.\n", name, claims.Subject)
}