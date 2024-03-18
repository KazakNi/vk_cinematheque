package api

import (
	"log"
	"net/http"
)

func SetRoutes() {

	mux := http.NewServeMux()

	signUpHandler := http.HandlerFunc(SignUp)
	loginHandler := http.HandlerFunc(SignIn)

	GetActorsHandler := http.HandlerFunc(GetListActors)
	createActorHanlder := http.HandlerFunc(CreateActor)
	deleteActorHandler := http.HandlerFunc(DeleteActor)
	putActorHandler := http.HandlerFunc(UpdateActor)

	createFilmHanlder := http.HandlerFunc(CreateFilm)
	deleteFilmHandler := http.HandlerFunc(DeleteFilm)
	putFilmHandler := http.HandlerFunc(UpdateFilm)
	getFilmsHandler := http.HandlerFunc(GetListFilms)

	// Swagger specification
	mux.HandleFunc("GET /redoc", ReDoc)
	mux.Handle("/films.yaml", http.FileServer(http.Dir("../api/static")))

	mux.Handle("POST /user", LogRequest(signUpHandler))
	mux.Handle("POST /user/login", LogRequest(loginHandler))

	mux.Handle("GET /films", LogRequest(AuthRequiredCheck(getFilmsHandler)))
	mux.Handle("POST /films", LogRequest(AuthRequiredCheck(IsAdminCheck(createFilmHanlder))))
	mux.Handle("DELETE /films/{id}", LogRequest(AuthRequiredCheck(IsAdminCheck(deleteFilmHandler))))
	mux.Handle("PUT /films/{id}", LogRequest(AuthRequiredCheck(IsAdminCheck(putFilmHandler))))

	mux.Handle("GET /actors", LogRequest(AuthRequiredCheck(GetActorsHandler)))
	mux.Handle("POST /actors", LogRequest(AuthRequiredCheck(IsAdminCheck(createActorHanlder))))
	mux.Handle("DELETE /actors/{id}", LogRequest(AuthRequiredCheck(IsAdminCheck(deleteActorHandler))))
	mux.Handle("PUT /actors/{id}", LogRequest(AuthRequiredCheck(IsAdminCheck(putActorHandler))))

	log.Println("Starting server")
	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		log.Fatal(err)
	}
}
