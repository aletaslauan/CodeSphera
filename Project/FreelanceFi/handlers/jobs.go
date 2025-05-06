package handlers

import (
	"encoding/json"
	//"fmt"              // Added for error formatting
	"freelancefi/db"   // Keep if needed (e.g. ListOpenJobs uses it directly)
	"freelancefi/services"
	"log"
	"net/http"
	"strconv"
	"time"
	//"github.com/jackc/pgx/v5/pgtype" // Import pgtype if service uses it directly
)

type JobsHandler struct {
	Service *services.JobService
}

// DTO for creating a job via API request (JSON)
type CreateJobRequest struct {
	CategoryID  int32    `json:"category_id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	BudgetMin   float64  `json:"budget_min"`
	BudgetMax   *float64 `json:"budget_max,omitempty"` // Use pointer for optional
	Deadline    *string  `json:"deadline,omitempty"` // format: YYYY-MM-DD, use pointer for optional
}

// Handler for POST /jobs (Accepts JSON)
func (h *JobsHandler) CreateJob(w http.ResponseWriter, r *http.Request) {
	var req CreateJobRequest
	// Expect JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("JSON decode error in CreateJob: %v", err)
		http.Error(w, "Bad request: Invalid JSON format", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// --- Safe User ID Retrieval ---
	userIDValue := r.Context().Value("userID") // Use your actual context key string
	if userIDValue == nil {
        log.Println("Error in CreateJob: userID not found in context")
		http.Error(w, "Unauthorized: Missing user information", http.StatusUnauthorized)
		return
	}
	userID, ok := userIDValue.(int)
	if !ok {
		log.Printf("Error in CreateJob: userID in context is not int: %T", userIDValue)
		http.Error(w, "Internal Server Error: Unexpected user ID type", http.StatusInternalServerError)
		return
	}
	// --- End Safe User ID retrieval ---

	// Parse optional deadline string from JSON request
	var deadlineTime *time.Time
	if req.Deadline != nil && *req.Deadline != "" {
		t, err := time.Parse("2006-01-02", *req.Deadline)
		if err != nil {
            log.Printf("Invalid deadline format in JSON request: %v", err)
			http.Error(w, "Bad request: Invalid deadline format (expected YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
		deadlineTime = &t
	}

	// Prepare input for the service layer
	jobInput := services.CreateJobInput{
		ClientID:    int32(userID),
		CategoryID:  req.CategoryID,
		Title:       req.Title,
		Description: req.Description,
		BudgetMin:   req.BudgetMin,
		BudgetMax:   req.BudgetMax,
		Deadline:    deadlineTime,
	}

	// Call the service
	job, err := h.Service.CreateJob(r.Context(), jobInput)
	if err != nil {
		log.Printf("Service error creating job (API): %v", err)
		// TODO: Inspect error from service for more specific client errors (4xx)
		http.Error(w, "Could not create job", http.StatusInternalServerError)
		return
	}

	// Respond with created job (JSON)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(job)
}

// *** NEW Handler for POST /jobsform (Accepts HTML Form Data) ***
func (h *JobsHandler) CreateJobFromForm(w http.ResponseWriter, r *http.Request) {
	
	// Parse the form data from the request
	if err := r.ParseForm(); err != nil {
		log.Printf("Form parse error in CreateJobFromForm: %v", err)
		http.Error(w, "Bad request: Could not parse form data", http.StatusBadRequest)
		return
	}

	// --- Safe User ID Retrieval ---
	userIDValue := r.Context().Value("userID") // Use your actual context key string
	if userIDValue == nil {
        log.Println("Error in CreateJobFromForm: userID not found in context")
		http.Error(w, "Unauthorized: Missing user information", http.StatusUnauthorized)
		return
	}
	userID, ok := userIDValue.(int)
	if !ok {
		log.Printf("Error in CreateJobFromForm: userID in context is not int: %T", userIDValue)
		http.Error(w, "Internal Server Error: Unexpected user ID type", http.StatusInternalServerError)
		return
	}
	// --- End Safe User ID retrieval ---

	// Extract and validate form values
	title := r.FormValue("title")
	description := r.FormValue("description")
	categoryIDStr := r.FormValue("category_id")
	budgetMinStr := r.FormValue("budget_min")
	budgetMaxStr := r.FormValue("budget_max") // Optional
	deadlineStr := r.FormValue("deadline")    // Optional

	// Basic validation (add more as needed)
	if title == "" || description == "" || categoryIDStr == "" || budgetMinStr == "" {
		// Consider using HTMX to return an error snippet if applicable
        log.Printf("Validation error in CreateJobFromForm: Missing required fields (title=%s, desc=%s, cat=%s, minbud=%s)", title, description, categoryIDStr, budgetMinStr)
		http.Error(w, "Bad request: Title, Description, Category, and Min Budget are required.", http.StatusBadRequest)
		return
	}

	// Convert values
	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 32)
	if err != nil {
		log.Printf("Form value conversion error (category_id) in CreateJobFromForm: %v", err)
		http.Error(w, "Bad request: Invalid category ID", http.StatusBadRequest)
		return
	}

	budgetMin, err := strconv.ParseFloat(budgetMinStr, 64)
	if err != nil {
		log.Printf("Form value conversion error (budget_min) in CreateJobFromForm: %v", err)
		http.Error(w, "Bad request: Invalid minimum budget", http.StatusBadRequest)
		return
	}

	var budgetMax *float64 // Use pointer for optional field
	if budgetMaxStr != "" {
		val, err := strconv.ParseFloat(budgetMaxStr, 64)
		if err != nil {
			log.Printf("Form value conversion error (budget_max) in CreateJobFromForm: %v", err)
			http.Error(w, "Bad request: Invalid maximum budget", http.StatusBadRequest)
			return
		}
		// Add validation: max >= min
		if val < budgetMin {
            log.Printf("Validation error in CreateJobFromForm: Max budget %f < Min budget %f", val, budgetMin)
			http.Error(w, "Bad request: Maximum budget cannot be less than minimum budget", http.StatusBadRequest)
			return
		}
		budgetMax = &val // Assign address of parsed float
	}

	var deadlineTime *time.Time // Use pointer for optional field
	if deadlineStr != "" {
		t, err := time.Parse("2006-01-02", deadlineStr) // HTML date input format
		if err != nil {
			log.Printf("Form value conversion error (deadline) in CreateJobFromForm: %v", err)
			http.Error(w, "Bad request: Invalid deadline format (use YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
		deadlineTime = &t
	}

	// Prepare input for the service layer (same input struct as the JSON handler)
	jobInput := services.CreateJobInput{
		ClientID:    int32(userID),
		CategoryID:  int32(categoryID),
		Title:       title,
		Description: description,
		BudgetMin:   budgetMin,
		BudgetMax:   budgetMax,
		Deadline:    deadlineTime,
	}

	// Call the service
	_, err = h.Service.CreateJob(r.Context(), jobInput) // Resulting job object isn't needed for redirect
	if err != nil {
		log.Printf("Service error creating job (Form): %v", err)
		// TODO: Inspect error type from service for more specific client errors
		http.Error(w, "Could not create job due to a server error.", http.StatusInternalServerError) // Redirect back to form with error?
		return
	}

	// Success: Redirect the user
	log.Printf("Job created successfully via form by user %d", userID)
	http.Redirect(w, r, "/client/dashboard", http.StatusSeeOther) // Or "/jobspage" or other relevant page
}


// Handler for GET /jobs (JSON API - Lists Open Jobs)
// *** Note: This handler currently accesses Service.DB directly. ***
// *** It should ideally call a method on h.Service instead. ***
// *** Example: jobs, err := h.Service.ListOpenJobs(ctx, listInput) ***
func (h *JobsHandler) ListOpenJobs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context() // Get context once
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, err := strconv.Atoi(limitStr)
	// Add max limit check for safety
	if err != nil || limit <= 0 || limit > 100 { // Example max limit
		limit = 10 // Default limit
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	// --- Refactor Suggestion: Call Service Method ---
	// Replace this direct DB call with a call to a service method:
	// listInput := services.ListJobsInput{ Limit: int32(limit), Offset: int32(offset) }
	// jobs, err := h.Service.ListOpenJobs(ctx, listInput) // Assuming service method exists
	// If keeping direct DB access for now:
	jobs, err := h.Service.DB.ListOpenJobs(ctx, db.ListOpenJobsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	// --- End Refactor Suggestion ---

	// Handle potential errors from DB/Service call
	// The service method (if used) should ideally handle ErrNoRows.
	// If calling DB directly, handle it here:
	if err != nil && err.Error() != "no rows in result set" { // Basic check, use errors.Is(err, pgx.ErrNoRows) if using pgx directly
		log.Printf("Error listing open jobs (API): %v", err)
		http.Error(w, "Could not list jobs", http.StatusInternalServerError)
		return
	}
	if jobs == nil { // Ensure empty JSON array [] instead of null
		jobs = []db.Job{}
	}

	// Respond with jobs list (JSON)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(jobs)
}