package api

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

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

type ResponseWriterWrapper struct {
	w          *http.ResponseWriter
	body       *bytes.Buffer
	statusCode *int
}

func (rww ResponseWriterWrapper) String() string {
	var buf bytes.Buffer

	buf.WriteString("Response:")

	buf.WriteString("Headers:")
	for k, v := range (*rww.w).Header() {
		buf.WriteString(fmt.Sprintf("%s: %v", k, v))
	}

	buf.WriteString(fmt.Sprintf(" Status Code: %d", *(rww.statusCode)))

	buf.WriteString("Body")
	buf.WriteString(rww.body.String())
	return buf.String()
}
func (rww ResponseWriterWrapper) Write(buf []byte) (int, error) {
	rww.body.Write(buf)
	return (*rww.w).Write(buf)
}

func (rww ResponseWriterWrapper) Header() http.Header {
	return (*rww.w).Header()

}

func (rww ResponseWriterWrapper) WriteHeader(statusCode int) {
	(*rww.statusCode) = statusCode
	(*rww.w).WriteHeader(statusCode)
}
func NewResponseWriterWrapper(w http.ResponseWriter) ResponseWriterWrapper {
	var buf bytes.Buffer
	var statusCode int = 200
	return ResponseWriterWrapper{
		w:          &w,
		body:       &buf,
		statusCode: &statusCode,
	}
}

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("income request: endpoint %s, method %s, rBody %s", r.URL, r.Method, r.Body)
		defer func() {
			rww := NewResponseWriterWrapper(w)
			log.Println(
				fmt.Sprintf(
					"[Request: %s] [Response: %s]",
					r, rww.String(),
				))
		}()
		next.ServeHTTP(w, r)

	})
}
