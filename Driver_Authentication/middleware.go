package Driver_Authentication

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nonAuthRoutes := []string{"/",  "/auth/users", "/auth/login", "/auth/signup", "/rosters", "/directions"}
		currentRoute := r.RequestURI
		requireAuth := true
		for _, route := range nonAuthRoutes {
			if route  == currentRoute {
				requireAuth = false
			}
		}

		// only roster get req is unauthenticated
		if r.Method != http.MethodGet && currentRoute == "/rosters" {
			requireAuth = true
		}

		if requireAuth {

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
			next.ServeHTTP(w, r)
		} else {
			next.ServeHTTP(w, r)
		}

	})
}