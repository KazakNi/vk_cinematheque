package api

import (
	"bytes"
	"cinematheque/internal/db"
	"cinematheque/internal/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

type TestEndpoint struct {
	name           string
	funcEndpoint   handler
	method         string
	endpoint       string
	body           []byte
	expectedStatus int
	expectedBody   []byte
}

func CreateNonAdminUser() db.User {
	var dbuser db.User
	dbuser.Id = 0
	dbuser.Email = "test@test.ru"
	dbuser.Password = "lolkek"
	dbuser.Is_admin = false
	return dbuser
}

func TearUpDB() {

	host, port, user, password, _, driver := utils.GetEnv()

	connUrl := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=test sslmode=disable",
		host, port, user, password)

	DB, err := sqlx.Connect(driver, connUrl)

	if err != nil {
		panic(fmt.Sprintf("%s, %s", err, connUrl))
	}

	err = DB.Ping()
	if err != nil {
		panic(err)
	}

	log.Println("Successfully connected to DB")
	migration_path := utils.GetTestMigrPathUp()
	// create users table and admin
	db.ExecMigration(DB, migration_path)
	db.DBConnection = DB
	db.ExecMigration(DB, migration_path)
}

func TearDown() {
	migration_path := utils.GetTestMigrPathDown()
	db.ExecMigration(db.DBConnection, migration_path)
	db.DBConnection.Close()
}

type handler func(http.ResponseWriter, *http.Request)

func authEndpointPing(method, endpoint string, body []byte, funcEndpoint handler) ([]byte, int) {

	req := httptest.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	funcEndpoint(w, req)

	res := w.Result()
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	if endpoint == "/user/login" {
		return nil, res.StatusCode
	}
	return b, res.StatusCode
}

func endpointPingAuthMiddleware(method, endpoint string, body []byte, funcEndpoint handler) ([]byte, int) {

	req := httptest.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(funcEndpoint)
	handlerToTest := AuthRequiredCheck(nextHandler)

	handlerToTest.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return b, res.StatusCode
}

func TestAuthSignInAndUp(t *testing.T) {

	TearUpDB()

	user := CreateNonAdminUser()
	userJSON := User{Email: user.Email,
		Password: user.Password,
	}
	userBody, err := json.Marshal(&userJSON)

	if err != nil {
		log.Println(err)
		return
	}

	var tests = []TestEndpoint{
		{
			name:           "signup",
			funcEndpoint:   SignUp,
			method:         http.MethodPost,
			endpoint:       "/user",
			body:           userBody,
			expectedBody:   []byte(fmt.Sprintf(`{"id":1}`)),
			expectedStatus: 201,
		},
		{
			name:           "signin",
			funcEndpoint:   SignIn,
			method:         http.MethodPost,
			endpoint:       "/user/login",
			body:           userBody,
			expectedBody:   nil,
			expectedStatus: 200,
		},
	}

	for _, tc := range tests {
		body, code := authEndpointPing(tc.method, tc.endpoint, tc.body, tc.funcEndpoint)
		assert.Equal(t, tc.expectedBody, body)
		assert.Equal(t, tc.expectedStatus, code)
	}
	TearDown()

}

func TestGetRequestsNoAuth(t *testing.T) {

	var tests = []TestEndpoint{
		{
			name:           "getActors",
			funcEndpoint:   GetListActors,
			method:         http.MethodGet,
			endpoint:       "/actors",
			body:           nil,
			expectedBody:   []byte("Auth required!"),
			expectedStatus: 401,
		},
		{
			name:           "getFilms",
			funcEndpoint:   GetListFilms,
			method:         http.MethodGet,
			endpoint:       "/films",
			body:           nil,
			expectedBody:   []byte("Auth required!"),
			expectedStatus: 401,
		},
		{
			name:           "postFilm",
			funcEndpoint:   CreateFilm,
			method:         http.MethodPost,
			endpoint:       "/films",
			body:           nil,
			expectedBody:   []byte("Auth required!"),
			expectedStatus: 401,
		},
		{
			name:           "postActor",
			funcEndpoint:   CreateActor,
			method:         http.MethodPost,
			endpoint:       "/actors",
			body:           nil,
			expectedBody:   []byte("Auth required!"),
			expectedStatus: 401,
		},
		{
			name:           "deleteActor",
			funcEndpoint:   DeleteActor,
			method:         http.MethodDelete,
			endpoint:       "/actors/1",
			body:           nil,
			expectedBody:   []byte("Auth required!"),
			expectedStatus: 401,
		},
	}
	for _, tc := range tests {
		body, code := endpointPingAuthMiddleware(tc.method, tc.endpoint, tc.body, tc.funcEndpoint)
		assert.Equal(t, tc.expectedBody, body)
		assert.Equal(t, tc.expectedStatus, code)
	}

}

func TestCreateRequestsAuth(t *testing.T) {
	TearUpDB()

	user := CreateNonAdminUser()
	userJSON := User{Email: user.Email,
		Password: user.Password,
	}
	userBody, err := json.Marshal(&userJSON)

	if err != nil {
		log.Println(err)
		return
	}

	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(userBody))
	w := httptest.NewRecorder()

	SignUp(w, req)

	req = httptest.NewRequest(http.MethodPost, "/user/login", bytes.NewBuffer(userBody))
	w = httptest.NewRecorder()

	SignIn(w, req)
	res := w.Result()
	defer res.Body.Close()

	auth_header := res.Header["Authorization"]

	tt, _ := time.Parse("2006-01-02", "1990-02-12")

	a := Actor{
		Name:       "Lolkek",
		Sex:        "male",
		Birth_date: BirthDate{tt},
	}

	actor, err := json.Marshal(&a)

	req = httptest.NewRequest(http.MethodPost, "/actors", bytes.NewBuffer(actor))
	w = httptest.NewRecorder()
	req.Header.Set("Authorization", strings.Join(auth_header, " "))

	nextHandler := http.HandlerFunc(CreateActor)
	handlerToTest := AuthRequiredCheck(nextHandler)

	handlerToTest.ServeHTTP(w, req)

	res = w.Result()
	defer res.Body.Close()

	id := &CreatedId{}
	d := json.NewDecoder(res.Body)

	err = d.Decode(id)
	if err != nil {
		log.Println(err)
		return
	}

	assert.Equal(t, 201, res.StatusCode)
	assert.Equal(t, 1, id.Id)

}
