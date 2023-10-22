package http

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"refactoring/internal/models"
	"reflect"
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
	w.Header().Set("Content-Type", "application/json")

	var user models.CreateUserRequest

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		err = json.NewEncoder(w).Encode(Error{ErrorMessage: err.Error()})
		if err != nil {
			h.logger.Printf("responseWriter error: %v", err)
		}
		return
	}
	if reflect.ValueOf(user).IsZero() {
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(Error{ErrorMessage: "empty_values"})
		if err != nil {
			h.logger.Printf("responseWriter error: %v", err)
		}
		return
	}

	createUser, err := h.client.CreateUser(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode(Error{ErrorMessage: err.Error()})
		if err != nil {
			h.logger.Printf("responseWriter error: %v", err)
		}
		return
	}
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(createUser)
	if err != nil {
		h.logger.Printf("responseWriter error: %v", err)
	}

	return

}

func (h *HttpRouter) getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")
	user, err := h.client.GetUser(id)

	if err != nil {

		if errors.Is(err, models.UserNotFound) {
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
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")
	var user models.CreateUserRequest

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		err = json.NewEncoder(w).Encode(Error{ErrorMessage: err.Error()})
		if err != nil {
			h.logger.Printf("responseWriter error: %v", err)
		}
		return
	}
	if user.DisplayName == "" {
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(Error{ErrorMessage: "empty_values"})
		if err != nil {
			h.logger.Printf("responseWriter error: %v", err)
		}
		return
	}

	updateUser, err := h.client.UpdateUser(id, user.DisplayName)
	if err != nil {
		if errors.Is(err, models.UserNotFound) {
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

	err = json.NewEncoder(w).Encode(updateUser)
	if err != nil {
		h.logger.Printf("responseWriter error: %v", err)
	}

	return

}

func (h *HttpRouter) deleteUser(w http.ResponseWriter, r *http.Request) {

}
