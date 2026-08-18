package services

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Handler struct {
	store *Store
}

type createServiceRequest struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Status  string `json:"status"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func NewHandler(store *Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req createServiceRequest
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON request body")
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Version = strings.TrimSpace(req.Version)
	req.Status = strings.TrimSpace(req.Status)
	if req.Name == "" || req.Version == "" || req.Status == "" {
		writeError(w, http.StatusBadRequest, "name, version, and status are required")
		return
	}

	ctx, cancel := contextWithTimeout(r, 3*time.Second)
	defer cancel()

	service, err := h.store.Create(ctx, CreateServiceParams{
		Name:    req.Name,
		Version: req.Version,
		Status:  req.Status,
	})
	if err != nil {
		log.Printf("create service failed: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to create service")
		return
	}

	writeJSON(w, http.StatusCreated, service)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := contextWithTimeout(r, 3*time.Second)
	defer cancel()

	services, err := h.store.List(ctx)
	if err != nil {
		log.Printf("list services failed: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to list services")
		return
	}

	if services == nil {
		services = []Service{}
	}

	writeJSON(w, http.StatusOK, services)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid service id")
		return
	}

	ctx, cancel := contextWithTimeout(r, 3*time.Second)
	defer cancel()

	service, err := h.store.Get(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		writeError(w, http.StatusNotFound, "service not found")
		return
	}
	if err != nil {
		log.Printf("get service failed: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to get service")
		return
	}

	writeJSON(w, http.StatusOK, service)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("write response failed: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}
