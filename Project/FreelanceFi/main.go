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
	"freelancefi/services"
	"freelancefi/templates"
	clienttpl "freelancefi/templates/client"
	freelencertpl "freelancefi/templates/freelancer"

	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5" 
	"golang.org/x/crypto/bcrypt"
	
)

// Context key for user ID (recommended practice).
// Define this at the package level so all handlers in this file can access it.
type contextKey string
const userIDKey contextKey = "userID"

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

	jobService := services.JobService{DB: queries}
	jobsHandler := handlers.JobsHandler{Service: &jobService}

	
	bidService := services.BidService{DB: queries}
	bidsHandler := handlers.BidsHandler{Service: &bidService}

	r.With(AuthMiddleware).Post("/jobs/{jobID}/bids", bidsHandler.PlaceBid)
	r.With(AuthMiddleware).Get("/jobs/{jobID}/bids", bidsHandler.ListBidsForJob)

	
	userHandler := handlers.UserHandler{Queries: queries}
	//jobsHandler := handlers.JobsHandler{Queries: queries}
	//bidsHandler := handlers.BidsHandler{Queries: queries}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		templates.LoginPage().Render(r.Context(), w)
	})
	r.Get("/register", func(w http.ResponseWriter, r *http.Request) {
		templates.RegisterPage().Render(r.Context(), w)
	})

	r.With(AuthMiddleware).Get("/home", homeHandler(queries))
	r.With(AuthMiddleware).Get("/profile", profileHandler(queries))
	r.With(AuthMiddleware).Get("/jobspage", jobspageHandler(queries))

	r.With(AuthMiddleware).Get("/createjobs", createJobPageHandler(queries))
	//r.Post("/jobsform", jobsHandler.CreateJobFromForm) 
	r.With(AuthMiddleware).Post("/jobsform", jobsHandler.CreateJobFromForm)

	r.With(AuthMiddleware).Get("/mywork", myworkHandler(queries))
	r.With(AuthMiddleware).Get("/finance", financeHandler(queries))

	r.With(AuthMiddleware).Get("/client/dashboard", clientDashboardHandler(queries))
	r.With(AuthMiddleware).Get("/freelancer/freelancer_dashboard", freelancerDashboardHandler(queries))

	r.With(AuthMiddleware).Post("/jobs", jobsHandler.CreateJob)
	r.With(AuthMiddleware).Get("/jobs", jobsHandler.ListOpenJobs)

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
			w.Header().Set("HX-Redirect", "/freelancer/freelancer_dashboard")
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
				Registration successful! Redirecting to I page...
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

func jobspageHandler(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int)
		user, err := queries.GetUser(r.Context(), int32(userID))
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		query := r.URL.Query()
		limit := 9 // default
		offset := 0

		if l := query.Get("limit"); l != "" {
			if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
				limit = parsed
			}
		}

		if o := query.Get("offset"); o != "" {
			if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
				offset = parsed
			}
		}

		totalJobs, err := queries.CountOpenJobs(r.Context())
		if err != nil {
			http.Error(w, "Could not count jobs", http.StatusInternalServerError)
			return
		}

		jobs, err := queries.ListOpenJobs(r.Context(), db.ListOpenJobsParams{
			Limit:  int32(limit),
			Offset: int32(offset),
		})
		if err != nil {
			http.Error(w, "Could not fetch jobs", http.StatusInternalServerError)
			return
		}

		prevOffset := offset - limit
		if prevOffset < 0 {
			prevOffset = 0
		}
		nextOffset := offset + limit

		currentPage := (offset / limit) + 1
		totalPages := (int(totalJobs) + limit - 1) / limit

		err = templates.JobsPage(user.Username, jobs, limit, offset, prevOffset, nextOffset, currentPage, totalPages).Render(r.Context(), w)
		if err != nil {
			log.Println("template render error:", err)
			http.Error(w, "Template error", http.StatusInternalServerError)
		}
	}
}

func createJobPageHandler(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int)
		user, err := queries.GetUser(r.Context(), int32(userID))
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		categories, err := queries.ListCategories(r.Context()) // your SQL query for categories
		if err != nil {
			http.Error(w, "Failed to load categories", http.StatusInternalServerError)
			return
		}

		var names []string
		for _, cat := range categories {
			names = append(names, cat.Name)
		}

		clienttpl.CreateJobPage(user.Username, categories).Render(r.Context(), w)
	}
}





func handleJobBidsPage(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// --- Get Job ID ---
		jobIDStr := chi.URLParam(r, "jobID")
		jobID64, err := strconv.ParseInt(jobIDStr, 10, 32)
		if err != nil {
			log.Printf("Invalid job ID format in URL: %s", jobIDStr)
			http.Error(w, "Invalid Job ID", http.StatusBadRequest)
			return
		}
		jobID := int32(jobID64)
		// --- Get Current User ---
		userIDValue := ctx.Value(userIDKey)
		if userIDValue == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		currentUserID, ok := userIDValue.(int)
		if !ok {
			log.Printf("Error: userID in context is not int: %T", userIDValue)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		currentUserID32 := int32(currentUserID)
		currentUser, err := queries.GetUser(ctx, currentUserID32)
		if err != nil {
			log.Printf("Error fetching user %d: %v", currentUserID, err)
			// Handle error appropriately (e.g., 500 or 404)
			http.Error(w, "Error fetching user", http.StatusInternalServerError)
			return
		}
		// --- Fetch Job ---
		job, err := queries.GetJobByID(ctx, jobID) // Use correct GetJobByID
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				http.NotFound(w, r)
				return
			}
			log.Printf("Error fetching job %d: %v", jobID, err)
			http.Error(w, "Failed to load job", http.StatusInternalServerError)
			return
		}
		// --- Fetch Bids ---
		bids, err := queries.ListBidsForJob(ctx, jobID)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			log.Printf("Error fetching bids %d: %v", jobID, err)
			http.Error(w, "Failed to load bids", http.StatusInternalServerError)
			return
		}
		if errors.Is(err, pgx.ErrNoRows) || bids == nil {
			bids = []db.Bid{}
		}
		// --- Render Template ---
		component := clienttpl.BidsPage(currentUser.Username, currentUserID32, job, bids)
		renderErr := component.Render(ctx, w)
		if renderErr != nil {
			log.Printf("Error rendering BidsPage %d: %v", jobID, renderErr)
		}
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
		clienttpl.ClientPage(user.Username, user.Role, "/client/profile", "Go to your profile").Render(r.Context(), w)

		//templates.ClientDashboard(user.Username).Render(r.Context(), w)
		//w.Header().Set("HX-Redirect", "/client/dashboard")
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
		freelencertpl.HomePage(user.Username, user.Role, "/profile", "Go to your profile").Render(r.Context(), w)

		//w.Header().Set("HX-Redirect", "/freelancer/dashboard")
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
