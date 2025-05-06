package handlers

import (
	"encoding/json"
	"freelancefi/services"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type BidsHandler struct {
	Service *services.BidService
}

type PlaceBidRequest struct {
	Amount      float64 `json:"amount"`
	CoverLetter string  `json:"cover_letter"`
}

func (h *BidsHandler) PlaceBid(w http.ResponseWriter, r *http.Request) {
	jobIDStr := chi.URLParam(r, "jobID")
	jobID, err := strconv.Atoi(jobIDStr)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	var req PlaceBidRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(int)

	bid, err := h.Service.PlaceBid(r.Context(), services.PlaceBidInput{
		JobID:        int32(jobID),
		FreelancerID: int32(userID),
		Amount:       req.Amount,
		CoverLetter:  req.CoverLetter,
	})
	if err != nil {
		log.Printf("place bid error: %v", err)
		http.Error(w, "Could not place bid", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(bid)
}

func (h *BidsHandler) ListBidsForJob(w http.ResponseWriter, r *http.Request) {
	jobIDStr := chi.URLParam(r, "jobID")
	jobID, err := strconv.Atoi(jobIDStr)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	bids, err := h.Service.DB.ListBidsForJob(r.Context(), int32(jobID))
	if err != nil {
		log.Printf("list bids error: %v", err)
		http.Error(w, "Could not fetch bids", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bids)
}

