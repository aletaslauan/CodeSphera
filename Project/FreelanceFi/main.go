package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"freelancefi/config"
	"freelancefi/db"
	"freelancefi/handlers"
	"freelancefi/templates"

	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func createSessionToken(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(72 * time.Hour).Unix(),
	})
	return token.SignedString(config.SecretKey)
}

func parseToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return config.SecretKey, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, err
	}
	return &claims, nil
}

func main() {
	config.InitConfig()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	pool := db.NewConnectionPool()
	queries := db.New(pool)

	userHandler := handlers.UserHandler{Queries: queries}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		templates.LoginPage().Render(r.Context(), w)
	})
	r.Get("/register", func(w http.ResponseWriter, r *http.Request) {
		templates.RegisterPage().Render(r.Context(), w)
	})

	r.With(AuthMiddleware).Get("/home", homeHandler(queries))
	r.With(AuthMiddleware).Get("/profile", profileHandler(queries))
	r.With(AuthMiddleware).Get("/jobs", jobsHandler(queries))
	r.With(AuthMiddleware).Get("/mywork", myworkHandler(queries))
	r.With(AuthMiddleware).Get("/finance", financeHandler(queries))

	r.With(AuthMiddleware).Get("/client/dashboard", clientDashboardHandler(queries))
	r.With(AuthMiddleware).Get("/freelancer/dashboard", freelancerDashboardHandler(queries))

	r.Post("/login", loginHandler(queries))
	r.Post("/register", registerHandler(queries))

	r.Get("/users", userHandler.GetUsers)
	r.Post("/users", userHandler.CreateUser)

	r.With(AuthMiddleware).Get("/protected", protectedHandler(queries))

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	fmt.Println("Server started at http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}

func loginHandler(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var username, password string
		contentType := r.Header.Get("Content-Type")

		if strings.HasPrefix(contentType, "application/json") {
			var creds handlers.CreateUserDTO
			if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
				http.Error(w, "Invalid JSON body", http.StatusBadRequest)
				return
			}
			username = creds.Username
			password = creds.Password
		} else {
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Invalid form", http.StatusBadRequest)
				return
			}
			username = r.FormValue("username")
			password = r.FormValue("password")
		}

		if username == "" || password == "" {
			http.Error(w, "Username and password are required", http.StatusBadRequest)
			return
		}

		user, err := queries.GetUserByUsername(r.Context(), username)
		if err != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`<div class="alert alert-danger">Invalid username or password. Please try again.</div>`))
			return
		}

		token, err := createSessionToken(int(user.ID))
		if err != nil {
			http.Error(w, "Could not create session", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    token,
			Path:     "/",
			Expires:  time.Now().Add(72 * time.Hour),
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})

		if user.Role == "client" {
			w.Header().Set("HX-Redirect", "/client/dashboard")
		} else {
			w.Header().Set("HX-Redirect", "/freelancer/dashboard")
		}
	}
}

func registerHandler(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form", http.StatusBadRequest)
			return
		}
		username := r.FormValue("username")
		password := r.FormValue("password")
		role := r.FormValue("role")

		if username == "" || password == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`<div class="alert alert-danger">Username and password are required.</div>`))
			return
		}

		if role != "client" && role != "freelancer" {
			role = "freelancer"
		}

		_, err := queries.GetUserByUsername(r.Context(), username)
		if err == nil {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(`<div class="alert alert-danger">Username already taken. Please choose another.</div>`))
			return
		} else if !errors.Is(err, sql.ErrNoRows) {
			log.Printf("error checking user: %v", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		_, err = queries.AddUser(r.Context(), db.AddUserParams{
			Username:     username,
			PasswordHash: string(hashedPassword),
			Role:         role,
		})
		if err != nil {
			log.Printf("error creating user: %v", err)
			http.Error(w, "Could not create user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`
			<div class="alert alert-success text-center" role="alert">
				Registration successful! Redirecting to login page...
				<div class="spinner-border spinner-border-sm ms-2" role="status"></div>
			</div>
			<script>
				setTimeout(function() {
					window.location.href = "/";
				}, 1000);
			</script>
		`))
	}
}

func profileHandler(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int)
		user, err := queries.GetUser(r.Context(), int32(userID))
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		templates.ProfilePage(user.Username).Render(r.Context(), w)
	}
}

func homeHandler(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int)
		user, err := queries.GetUser(r.Context(), int32(userID))
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		templates.HomePage(user.Username, user.Role, "/profile", "Go to your profile").Render(r.Context(), w)
	}
}

func jobsHandler(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int)
		user, err := queries.GetUser(r.Context(), int32(userID))
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		templates.JobsPage(user.Username).Render(r.Context(), w)
	}
}

func myworkHandler(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int)
		user, err := queries.GetUser(r.Context(), int32(userID))
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		templates.MyWorkPage(user.Username).Render(r.Context(), w)
	}
}

func financeHandler(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int)
		user, err := queries.GetUser(r.Context(), int32(userID))
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		templates.FinancePage(user.Username).Render(r.Context(), w)
	}
}

func clientDashboardHandler(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int)
		user, err := queries.GetUser(r.Context(), int32(userID))
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		if user.Role != "client" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		//templates.ClientDashboard(user.Username).Render(r.Context(), w)
		w.Header().Set("HX-Redirect", "/client/dashboard")
	}
}

func freelancerDashboardHandler(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int)
		user, err := queries.GetUser(r.Context(), int32(userID))
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		if user.Role != "freelancer" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		w.Header().Set("HX-Redirect", "/freelancer/dashboard")
		// 	templates.FreelancerDashboard(user.Username).Render(r.Context(), w)
	}
}

func protectedHandler(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int)
		user, err := queries.GetUser(r.Context(), int32(userID))
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		claims, err := parseToken(cookie.Value)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		userID := int((*claims)["userID"].(float64))
		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
