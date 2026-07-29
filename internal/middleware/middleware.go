package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/HammerBone/project-run/internal/util"
)

type authContextKey struct{}

var Authkey = authContextKey{}

func GetClaims(ctx context.Context) (*util.UserSecretClaim, bool) {
	claims, ok := ctx.Value(Authkey).(*util.UserSecretClaim)
	return claims, ok
}

func AuthMiddleware(jwtGenerator *util.JWTGenerator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "authorization header is missing", http.StatusUnauthorized)
				return
			}

			fields := strings.Fields(authHeader)
			if len(fields) != 2 || fields[0] != "Bearer" {
				http.Error(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}

			token := fields[1]
			claims, err := jwtGenerator.VerifyToken(token)
			if err != nil {
				http.Error(w, "invalid claim", http.StatusUnauthorized)
				log.Printf("error: %+v", err)
				return
			}
			log.Println("[middleware][VerifyToken] claims: ", claims)

			ctx := context.WithValue(r.Context(), authContextKey{}, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
