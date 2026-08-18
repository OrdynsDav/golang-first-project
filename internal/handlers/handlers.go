package handlers

import (
	"encoding/json"
	"net/http"
	"restapi-tasks/internal/database"
	"restapi-tasks/internal/helpers"
	"restapi-tasks/internal/models"
	"strconv"
	"strings"
)

type Handlers struct {
	store *database.TaskStore
}

func NewHandlerStore(store *database.TaskStore) *Handlers {
	return &Handlers{store: store}
}

func (h *Handlers) GetAllTasks(w http.ResponseWriter, request *http.Request) {
	tasks, err := h.store.GetAll()

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Cannot get tasks")
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, tasks)
}

func (h *Handlers) GetTask(w http.ResponseWriter, request *http.Request) {
	pathParts := strings.Split(strings.TrimPrefix(request.URL.Path, "/tasks/"), "/")

	idString := pathParts[0]

	id, err := strconv.Atoi(idString)

	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	task, err := h.store.GetById(id)

	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Cannot get task")
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, task)
}

func (h *Handlers) CreateTask(w http.ResponseWriter, request *http.Request) {
	var input models.CreateTaskInput

	err := json.NewDecoder(request.Body).Decode(&input)

	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Incorrectly sent data")
		return
	}

	if strings.TrimSpace(input.Title) == "" {
		helpers.RespondWithError(w, http.StatusBadRequest, "Task title must be not empty")
		return
	}

	task, err := h.store.Create(input)

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.RespondWithJSON(w, http.StatusCreated, task)
}

func (h *Handlers) UpdateTask(w http.ResponseWriter, request *http.Request) {
	pathParts := strings.Split(strings.TrimPrefix(request.URL.Path, "/tasks/"), "/")

	idString := pathParts[0]

	id, err := strconv.Atoi(idString)

	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var input models.UpdateTaskInput

	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Incorrectly sent data")
		return
	}

	if input.Title != nil && strings.TrimSpace(*input.Title) == "" {
		helpers.RespondWithError(w, http.StatusBadRequest, "Task title must be not empty")
		return
	}

	task, err := h.store.Update(id, input)

	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			helpers.RespondWithError(w, http.StatusNotFound, err.Error())
		} else {
			helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, task)
}

func (h *Handlers) DeleteTask(w http.ResponseWriter, request *http.Request) {
	pathParts := strings.Split(strings.TrimPrefix(request.URL.Path, "/tasks/"), "/")

	idString := pathParts[0]

	id, err := strconv.Atoi(idString)

	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	err = h.store.Delete(id)

	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			helpers.RespondWithError(w, http.StatusNotFound, err.Error())
		} else {
			helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "success"})
}
