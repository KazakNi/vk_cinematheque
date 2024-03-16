package api

import (
	"log"
	"net/http"
)

func SetRoutes() {

	mux := http.NewServeMux()

	signupHandler := http.HandlerFunc(SignUp)

	mux.Handle("POST /user", CheckRBody(signupHandler))
	mux.HandleFunc("POST /user/login", SignIn)

	log.Println("Starting server")
	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		log.Fatal(err)
	}
}
