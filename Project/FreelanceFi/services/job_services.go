package services

import (
	"context"
	"freelancefi/db"
	"time"
	"fmt" 
	"log"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
)

type JobService struct {
	DB *db.Queries
}

type CreateJobInput struct {
	ClientID    int32
	CategoryID  int32
	Title       string
	Description string
	BudgetMin   float64
	BudgetMax   *float64
	Deadline    *time.Time
}

func (s *JobService) CreateJob(ctx context.Context, input CreateJobInput) (db.Job, error) {
	// Convert float64 → pgtype.Numeric
	budgetMin := pgtype.Numeric{}
if err := budgetMin.Scan(strconv.FormatFloat(input.BudgetMin, 'f', 2, 64)); err != nil {
	log.Printf("budgetMin scan failed: %v", err)
	return db.Job{}, fmt.Errorf("invalid budget_min: %w", err)
}


	// Optional max
	budgetMax := pgtype.Numeric{}
if input.BudgetMax != nil {
	strVal := strconv.FormatFloat(*input.BudgetMax, 'f', 2, 64)
	if err := budgetMax.Scan(strVal); err != nil {
		return db.Job{}, fmt.Errorf("invalid budget_max: %w", err)
	}
} else {
	budgetMax.Valid = false
}

	// Optional deadline
	deadline := pgtype.Date{}
	if input.Deadline != nil {
		if err := deadline.Scan(*input.Deadline); err != nil {
			return db.Job{}, fmt.Errorf("invalid deadline: %w", err)
		}
	} else {
		deadline.Valid = false
	}

	// Create job
	return s.DB.CreateJob(ctx, db.CreateJobParams{
		ClientID:    input.ClientID,
		CategoryID:  input.CategoryID,
		Title:       input.Title,
		Description: input.Description,
		BudgetMin:   budgetMin,
		BudgetMax:   budgetMax,
		Deadline:    deadline,
	})
}

