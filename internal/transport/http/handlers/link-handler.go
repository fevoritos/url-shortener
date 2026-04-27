package handlers

import (
	"encoding/json"
	"net/http"
	linkusecase "url-shortener/internal/usecase/link"
)

type LinkHandler struct {
	usecase linkusecase.Usecase
}

func NewLinkHandler(usecase linkusecase.Usecase) *LinkHandler {
	return &LinkHandler{usecase: usecase}
}

func (h *LinkHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req linkDTO
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	created, err := h.usecase.Create(r.Context(), req.Url)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
	}

	writeJSON(w, http.StatusCreated, created)
}

func (h *LinkHandler) Goto(w http.ResponseWriter, r *http.Request) {
	hash := r.PathValue("hash")
	link, err := h.usecase.GetByHash(r.Context(), hash)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
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

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{
		"error": err.Error(),
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(payload)
}
