package task

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func writeJson(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	writeJson(w, statusCode, map[string]string{"error": message})
}

func readJson(r *http.Request, target interface{}) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(target)
}

func (h *Handler) HandleCreate(w http.ResponseWriter, r *http.Request) {

	var req CreateRequest
	if err := readJson(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	task, err := h.service.CreateTask(r.Context(), req)

	if err != nil {
		if errors.Is(err, ErrEmptyTitle) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	writeJson(w, http.StatusCreated, task)
}

func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "Invalid task id")
		return
	}

	task, err := h.service.GetTask(r.Context(), id)

	if err != nil {
		if errors.Is(err, ErrNotFound) {
			writeError(w, http.StatusNotFound, "Task Not Found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	writeJson(w, http.StatusOK, task)
}

func (h *Handler) HandleList(w http.ResponseWriter, r *http.Request) {
	task, err := h.service.ListTask(r.Context())

	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	writeJson(w, http.StatusOK, task)
}

func (h *Handler) HandleUpdate(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "Invalid task id")
		return
	}

	var req UpdateRequest

	if err := readJson(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	task, err := h.service.UpdateTask(r.Context(), id, req)

	if err != nil {
		if errors.Is(err, ErrNotFound) {
			writeError(w, http.StatusNotFound, "task not found")
			return
		}

		if errors.Is(err, ErrEmptyTitle) || errors.Is(err, ErrInvalidStatus) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	writeJson(w, http.StatusOK, task)
}

func (h *Handler) HandleDelete(w http.ResponseWriter, r *http.Request) {

	idstr := r.PathValue("id")

	id, err := strconv.ParseInt(idstr, 10, 64)

	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "Invalid status id")
		return
	}

	err = h.service.DeleteTask(r.Context(), id)

	if err != nil {
		if errors.Is(err, ErrNotFound) {
			writeError(w, http.StatusNotFound, "Task not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
