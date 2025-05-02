// handlers/user.go
package handlers

import (
	"encoding/json"
	"freelancefi/db"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type CreateUserDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type UserDTO struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

type UserHandler struct {
	Queries *db.Queries
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.Queries.SelectUsers(r.Context())
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	usersDTO := make([]UserDTO, len(users))
	for i, u := range users {
		usersDTO[i] = UserDTO{
			Username: u.Username,
			Role:     u.Role, 
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usersDTO)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var decodedBody CreateUserDTO
	err := json.NewDecoder(r.Body).Decode(&decodedBody)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if decodedBody.Username == "" || decodedBody.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}
	if decodedBody.Role != "client" && decodedBody.Role != "freelancer" {
		decodedBody.Role = "freelancer" // Default to freelancer if invalid/missing
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(decodedBody.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	user, err := h.Queries.AddUser(r.Context(), db.AddUserParams{
		Username:     decodedBody.Username,
		PasswordHash: string(hashedPassword),
		Role:         decodedBody.Role, //Save role into DB
	})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(UserDTO{
		Username: user.Username,
		Role:     user.Role,
	})
}
