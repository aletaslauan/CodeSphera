package services

import (
	"context"
	"freelancefi/db"
	"strconv"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type BidService struct {
	DB *db.Queries
}

type PlaceBidInput struct {
	FreelancerID int32
	JobID        int32
	Amount       float64
	CoverLetter  string
}

func (s *BidService) PlaceBid(ctx context.Context, input PlaceBidInput) (db.Bid, error) {
	// Convert amount (float64) → pgtype.Numeric
	amount := pgtype.Numeric{}
	if err := amount.Scan(strconv.FormatFloat(input.Amount, 'f', 2, 64)); err != nil {
		return db.Bid{}, fmt.Errorf("invalid amount: %w", err)
	}

	// Convert cover letter (string) → pgtype.Text
	coverLetter := pgtype.Text{}
	if err := coverLetter.Scan(input.CoverLetter); err != nil {
		return db.Bid{}, fmt.Errorf("invalid cover letter: %w", err)
	}

	// Build params and call PlaceBid
	params := db.PlaceBidParams{
		JobID:        input.JobID,
		FreelancerID: input.FreelancerID,
		Amount:       amount,
		CoverLetter:  coverLetter,
	}

	return s.DB.PlaceBid(ctx, params)
}

