package main

import (
	"curs1_boilerplate/db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// ROUTING
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	// connect to the db
	pool := db.NewConnectionPool()
	queries := db.New(pool)

	// just a test route
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("server is up"))
	})

	// get all users
	type UserDTO struct {
		Username string `json:"user"`
	}
	r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		users, err := queries.SelectUsers(r.Context())
		if err != nil {
			log.Printf("select users error: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		usersDTO := make([]UserDTO, len(users))
		for i, u := range users {
			usersDTO[i] = UserDTO{Username: u.Username}
		}

		json.NewEncoder(w).Encode(usersDTO)
	})

	// add users
	type CreateUserDTO struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	r.Post("/users", func(w http.ResponseWriter, r *http.Request) {
		// get the username and password from the request body

		var decodedBody CreateUserDTO 
		err := json.NewDecoder(r.Body).Decode(&decodedBody)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		err = queries.AddUser(r.Context(), db.AddUserParams{
			Username: decodedBody.Username,
			Password: decodedBody.Password,
		})
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	})

	// START THE SERVER
	fmt.Printf("Server stared at http://localhost%s\n", PORT)
	log.Fatal(http.ListenAndServe(PORT, r))
}
