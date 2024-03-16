package api

import (
	"cinematheque/internal/db"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
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

func SignUp(w http.ResponseWriter, r *http.Request) {
	user := &User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		log.Printf("Error while %s endpoint response body parsing: %s", r.URL, err)
	}
	w.Header().Set("Content-Type", "application/json")
	userExists, err := user.IsUserExists(db.DBConnection)

	if err != nil {
		log.Printf("Error while checking user existence")
	}

	if userExists {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("409 - User already exists"))
		return
	}

	user.HashPassword(user.Password)
	userId, err := user.CreateUser(db.DBConnection)

	if err != nil {
		log.Printf("Error while creating user: %s", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	resp, err := json.Marshal(UserId{Id: userId})
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		return
	}
	w.Write(resp)
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	user := &User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		log.Printf("Error while %s endpoint response body parsing: %s", r.URL, err)
	}
	w.Header().Set("Content-Type", "application/json")

	db_user, err := user.GetUserByEmail(user.Email, db.DBConnection)

	if err != nil {
		log.Printf("Error while user querying: %s", err)
	}

	if db_user.Email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - Invalid credentials"))
		log.Println("Wrong email")
		return
	}

	if user.CheckPassword(user.Password, db_user.Password) != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - Invalid credentials"))
		log.Println("Wrong password")
		return
	}
	is_admin := db_user.Is_admin
	claims := CustomClaims{is_admin, jwt.RegisteredClaims{ID: strconv.Itoa(db_user.Id), ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour))}}
	SetToken(w, r, claims)

}

func SetToken(w http.ResponseWriter, r *http.Request, claims CustomClaims) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token_string, err := token.SignedString([]byte(load_secret()))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - Invalid token"))
	}
	var bearer = "Bearer " + token_string
	r.Header.Set("Authorization", bearer)

	w.Header().Set("Content-Type", "application/json")

	token_response := make(map[string]string)
	token_response["token"] = token_string
	res, err := json.Marshal(token_response)
	if err != nil {
		log.Println("Error while marshalling token")

	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}

func load_secret() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	parent := filepath.Dir(wd)

	err = godotenv.Load(parent + "./.env")
	if err != nil {
		panic("Error loading .env file")
	}
	return os.Getenv("TOKEN_SECRET")
}
