package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	linkdomain "url-shortener/internal/domain"
	linkusecase "url-shortener/internal/usecase/link"
)

type LinkHandler struct {
	usecase linkusecase.Usecase
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewLinkHandler(usecase linkusecase.Usecase) *LinkHandler {
	return &LinkHandler{usecase: usecase}
}

func (h *LinkHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req linkDTO
	if err := decodeJSON(r, &req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	created, err := h.usecase.Create(r.Context(), req.Url)
	if err != nil {
		if errors.Is(err, linkusecase.ErrInvalidURL) {
			h.writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		slog.Error("failed to create link", "err", err)
		h.writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.writeJSON(w, http.StatusCreated, created)
}

func (h *LinkHandler) Goto(w http.ResponseWriter, r *http.Request) {
	hash := r.PathValue("hash")
	link, err := h.usecase.GetByHash(r.Context(), hash)
	if err != nil {
		if errors.Is(err, linkdomain.ErrNotFound) {
			h.writeError(w, http.StatusNotFound, "link not found")
			return
		}

		slog.Error("failed to get link", "hash", hash, "err", err)
		h.writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	http.Redirect(w, r, link.Url, http.StatusTemporaryRedirect)
}

func decodeJSON(r *http.Request, dst any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return err
	}

	return nil
}

func (h *LinkHandler) writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: msg})
}

func (h *LinkHandler) writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(payload)
}
