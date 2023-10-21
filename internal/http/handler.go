package http

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"refactoring/internal/models"
	"refactoring/internal/service/cache"
)

type Error struct {
	ErrorMessage string `json:"error_message"`
}

type Response struct {
	Users        *models.UserStore
	User         *models.User
	ErrorMessage *Error
}

func (h *HttpRouter) searchUsers(w http.ResponseWriter, r *http.Request) {

}

func (h *HttpRouter) createUser(w http.ResponseWriter, r *http.Request) {

}

func (h *HttpRouter) getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")
	user, err := h.client.GetUser(id)
	if err != nil {
		if errors.Is(err, cache.UserNotFound) {
			w.WriteHeader(http.StatusNotFound)
			err = json.NewEncoder(w).Encode(Error{ErrorMessage: err.Error()})
			if err != nil {
				h.logger.Printf("responseWriter error: %v", err)
			}
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(Error{ErrorMessage: err.Error()})
			if err != nil {
				h.logger.Printf("responseWriter error: %v", err)
			}
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		h.logger.Printf("responseWriter error: %v", err)
	}
	return
}

func (h *HttpRouter) updateUser(w http.ResponseWriter, r *http.Request) {

}

func (h *HttpRouter) deleteUser(w http.ResponseWriter, r *http.Request) {

}
