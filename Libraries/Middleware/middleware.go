package Middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"net/http"
)

var jwtKey = []byte("my_secret_key")

type Claims struct {
	Username string `json:"Username"`
	jwt.StandardClaims
}


func AuthMiddleware(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		c, err := r.Cookie("token")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		jwtTokenString := c.Value
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(jwtTokenString, claims, func(token *jwt.Token)(interface{}, error){
			return jwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		context.Set(r, "driverName", claims.Username)
		next.ServeHTTP(w, r)


	})
}