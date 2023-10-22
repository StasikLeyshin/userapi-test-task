package http

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"refactoring/internal/models"
	"reflect"
)

type ErrorResponse struct {
	ErrorMessage string `json:"error_message"`
}

func newError(ErrorMessage string) *ErrorResponse {
	return &ErrorResponse{
		ErrorMessage: ErrorMessage,
	}
}

func (h *HttpRouter) errInvalidRequest(err error) *ErrorResponse {

	switch {

	case errors.Is(err, models.UserNotFound):
		return newError(models.UserNotFound.Error())

	case errors.Is(err, models.EmptyValues):
		return newError(models.EmptyValues.Error())

	case h.debug:
		return newError(err.Error())

	default:
		return newError("unknown_error")
	}
}

func (h *HttpRouter) writeErrorResponse(w http.ResponseWriter, status int, errorResponse *ErrorResponse) {
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(errorResponse)
	if err != nil {
		h.logger.Printf("responseWriter error: %v", err)
	}
}

func (h *HttpRouter) searchUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users, err := h.client.SearchUsers()

	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, h.errInvalidRequest(err))
		return
	}

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		h.logger.Printf("responseWriter error: %v", err)
	}

	return

}

func (h *HttpRouter) createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user models.CreateUserRequest

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, h.errInvalidRequest(err))
		return
	}

	if reflect.ValueOf(user).IsZero() {
		h.writeErrorResponse(w, http.StatusBadRequest, h.errInvalidRequest(models.EmptyValues))
		return
	}

	createUser, err := h.client.CreateUser(&user)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, h.errInvalidRequest(err))
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

	switch {
	case errors.Is(err, models.UserNotFound):
		h.writeErrorResponse(w, http.StatusNotFound, h.errInvalidRequest(models.UserNotFound))
		return

	case err != nil:
		h.writeErrorResponse(w, http.StatusInternalServerError, h.errInvalidRequest(err))
		return
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
		h.writeErrorResponse(w, http.StatusBadRequest, h.errInvalidRequest(err))
		return
	}

	if user.DisplayName == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, h.errInvalidRequest(models.EmptyValues))
		return
	}

	updateUser, err := h.client.UpdateUser(id, user.DisplayName)
	switch {
	case errors.Is(err, models.UserNotFound):
		h.writeErrorResponse(w, http.StatusNotFound, h.errInvalidRequest(models.UserNotFound))
		return

	case err != nil:
		h.writeErrorResponse(w, http.StatusInternalServerError, h.errInvalidRequest(err))
		return
	}

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(updateUser)
	if err != nil {
		h.logger.Printf("responseWriter error: %v", err)
	}

	return

}

func (h *HttpRouter) deleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")

	err := h.client.DeleteUser(id)
	switch {
	case errors.Is(err, models.UserNotFound):
		h.writeErrorResponse(w, http.StatusNotFound, h.errInvalidRequest(models.UserNotFound))
		return

	case err != nil:
		h.writeErrorResponse(w, http.StatusNotFound, h.errInvalidRequest(err))
		return

	}

	w.WriteHeader(http.StatusNoContent)

	return
}
