package api

import (
	"log"
	"net/http"
)

func SetRoutes() {

	mux := http.NewServeMux()

	createActorHanlder := http.HandlerFunc(CreateActor)
	deleteActorHandler := http.HandlerFunc(DeleteActor)
	putActorHandler := http.HandlerFunc(UpdateActor)

	createFilmHanlder := http.HandlerFunc(CreateFilm)
	deleteFilmHandler := http.HandlerFunc(DeleteFilm)
	putFilmHandler := http.HandlerFunc(UpdateFilm)

	// Swagger specification
	mux.HandleFunc("GET /redoc", ReDoc)
	mux.Handle("/films.yaml", http.FileServer(http.Dir("../api/static")))

	mux.HandleFunc("POST /user", SignUp)
	mux.HandleFunc("POST /user/login", SignIn)

	mux.Handle("POST /films", AuthRequiredCheck(IsAdminCheck(createFilmHanlder)))
	mux.Handle("DELETE /films/{id}", AuthRequiredCheck(IsAdminCheck(deleteFilmHandler)))
	mux.Handle("PUT /films/{id}", AuthRequiredCheck(IsAdminCheck(putFilmHandler)))

	mux.Handle("POST /actors", AuthRequiredCheck(IsAdminCheck(createActorHanlder)))
	mux.Handle("DELETE /actors/{id}", AuthRequiredCheck(IsAdminCheck(deleteActorHandler)))
	mux.Handle("PUT /actors/{id}", AuthRequiredCheck(IsAdminCheck(putActorHandler)))

	log.Println("Starting server")
	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		log.Fatal(err)
	}
}
