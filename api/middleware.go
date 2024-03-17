package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

func CheckRBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Body == http.NoBody {
			http.Error(w, "Please send a request body", 400)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func AuthRequiredCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		AuthHeader := r.Header.Get("Authorization")

		if len(AuthHeader) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token := strings.Fields(AuthHeader)[1]
		err := ParseToken(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(fmt.Sprintf("%s", err)))
			log.Println(err)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func ParseToken(mytoken string) error {

	token, err := jwt.ParseWithClaims(mytoken, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(load_secret()), nil
	})
	if err != nil {
		log.Println("Error during token parsing", err)
		return errors.New("invalid token")
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		if claims.RegisteredClaims.ExpiresAt.Unix() < time.Now().Unix() {
			return errors.New("token is expired")
		}
	} else {
		return errors.New("invalid token")
	}
	return nil

}

func IsAdminCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		AuthHeader := r.Header.Get("Authorization")

		header_token := strings.Fields(AuthHeader)[1]
		token, _ := jwt.ParseWithClaims(header_token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(load_secret()), nil
		})

		claims, _ := token.Claims.(*CustomClaims)

		if !claims.Is_admin {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Forbidden"))
			log.Printf("User ID %s is not admin\n", claims.ID)
			return
		}

		next.ServeHTTP(w, r)
	})
}
