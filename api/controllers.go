package api

import (
	"cinematheque/internal/db"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	user := &User{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	err := d.Decode(user)
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
		log.Printf("Error while creating an user: %s", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	resp, err := json.Marshal(CreatedId{Id: userId})
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

func ReDoc(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("../api/static/redoc.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(w, nil)
}

func GetListActors(w http.ResponseWriter, r *http.Request) {
	var a Actors
	actors := []Actors{}
	actors, err := a.GetListActors(db.DBConnection)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error while extracting films list: %s", err)
		return
	}
	b, _ := json.Marshal(actors)
	w.Write(b)
}

func CreateActor(w http.ResponseWriter, r *http.Request) {
	actor := &Actor{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	err := d.Decode(actor)
	if err != nil {
		log.Printf("Error while %s endpoint response body parsing: %s", r.URL, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	actor_id, err := actor.Create(db.DBConnection)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error while creating an actor: %s", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	resp, err := json.Marshal(CreatedId{Id: actor_id})

	if err != nil {
		log.Printf("Error happened in JSON marshal. Err: %s", err)
		return
	}
	w.Write(resp)

}

func DeleteActor(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/actors/"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Undefined id request Err: %s", err)
		return
	}
	var actor Actor
	actor, err = actor.GetActorById(id, db.DBConnection)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(err)
		return
	}
	err = actor.Delete(id, db.DBConnection)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error while delete actor operation: %s", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func UpdateActor(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/actors/"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Undefined id request Err: %s", err)
		return
	}

	actor := &Actor{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	err = d.Decode(actor)

	if err != nil {
		log.Printf("Error while %s endpoint response body parsing: %s", r.URL, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = actor.GetActorById(id, db.DBConnection)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(err)
		return
	}

	err = actor.Update(id, db.DBConnection)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error while update actor operation: %s", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func CreateFilm(w http.ResponseWriter, r *http.Request) {
	film := &PostFilm{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	err := d.Decode(film)

	if err != nil {
		log.Printf("Error while %s endpoint response body parsing: %s", r.URL, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var f = Film{Name: film.Name,
		Description:  film.Description,
		Release_date: film.Release_date,
		Rating:       film.Rating}

	film_id, err := f.Create(db.DBConnection)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error while creating an film: %s", err)
		return
	}

	err = f.InsertCast(film_id, film.Actors_list, db.DBConnection)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error while inserting actors: %s", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	resp, err := json.Marshal(CreatedId{Id: film_id})

	if err != nil {
		log.Printf("Error happened in JSON marshal. Err: %s", err)
		return
	}
	w.Write(resp)

}

func DeleteFilm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/films/"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Undefined id request Err: %s", err)
		return
	}
	var film Film
	film, err = film.GetFilmById(id, db.DBConnection)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(err)
		return
	}
	err = film.Delete(id, db.DBConnection)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error while delete actor operation: %s", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func UpdateFilm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/films/"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Undefined id request Err: %s", err)
		return
	}

	film := &PostFilm{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	err = d.Decode(film)

	if err != nil {
		log.Printf("Error while %s endpoint response body parsing: %s", r.URL, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var f = Film{Name: film.Name,
		Description:  film.Description,
		Release_date: film.Release_date,
		Rating:       film.Rating}

	_, err = f.GetFilmById(id, db.DBConnection)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(err)
		return
	}

	err = f.Update(id, film.Actors_list, db.DBConnection)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error while update actor operation: %s", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func GetListFilms(w http.ResponseWriter, r *http.Request) {
	sort_param := r.URL.Query().Get("sort")
	searchByFilmParam := r.URL.Query().Get("search_by_movieName")
	searchByActorParam := r.URL.Query().Get("search_by_actorName")

	var f Films
	films := []Films{}
	films, err := f.GetListFilms(db.DBConnection, sort_param, searchByFilmParam, searchByActorParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error while extracting films list: %s", err)
		return
	}
	b, _ := json.Marshal(films)
	w.Write(b)
}
