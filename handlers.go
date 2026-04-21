package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	DB *pgxpool.Pool
}

func (s *Server) IngestEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var req EventRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "invalid JSON payload", http.StatusBadRequest)
		return
	}

	if req.EventType == "" || req.Source == "" || len(req.Payload) == 0 {
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	_, err := s.DB.Exec(
		ctx,
		`INSERT INTO analytics_events (event_type, source, payload)
		 VALUES ($1, $2, $3)`,
		req.EventType,
		req.Source,
		req.Payload,
	)

	if err != nil {
		http.Error(w, "failed to persist event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
