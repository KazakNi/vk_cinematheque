package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetListActorsNotAuth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/actors", nil)
	w := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(GetListActors)
	handlerToTest := AuthRequiredCheck(nextHandler)

	handlerToTest.ServeHTTP(w, req)

	res := w.Result()

	defer res.Body.Close()

	assert.Equal(t, 401, res.StatusCode)
}

func TestGetListFilmsNoAuth(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "/films", nil)
	w := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(GetListFilms)
	handlerToTest := AuthRequiredCheck(nextHandler)

	handlerToTest.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, 401, res.StatusCode)
}
