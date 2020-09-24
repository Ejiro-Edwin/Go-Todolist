package v1

import (
	"encoding/json"
	"github.com/ejiro-edwin/todolist/internal/api/utils"
	"github.com/ejiro-edwin/todolist/internal/database"
	"github.com/ejiro-edwin/todolist/internal/model"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

//TodoAPI - provides REST for Users
type TodoAPI struct {
	DB database.Database
}

func SetTodoAPI(db database.Database, router *mux.Router) {
	api := TodoAPI{
		DB:db,
	}

	apis := []API{
		NewAPI("/users/{userID}/todos", "POST", api.Create),
		NewAPI("/users/{userID}/todos", "GET", api.List),
		NewAPI("/users/{userID}/todos/{todoID}", "GET", api.Get),
		NewAPI("/users/{userID}/todos/{todoID}", "PATCH", api.Update),
		NewAPI("/users/{userID}/todos/{todoID}", "DELETE", api.Delete),
	}

	for _, api := range apis {
		router.HandleFunc(api.Path, api.Func).Methods(api.Method)
	}
}

//POST - /users/{userID}/todos
func (api *TodoAPI) Create(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "todo.go -> Create()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["userID"])

	logger = logger.WithFields(logrus.Fields{
		"userID":    userID,
	})

	//Decode parameters
	var todo model.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		logger.WithError(err).Warn("could not decode parameters")
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	todo.UserID = &userID

	if err := todo.VerifyCreate(); err != nil {
		logger.WithError(err).Warn("Not all fields found")
		utils.WriteError(w, http.StatusBadRequest, "Not all fields found", map[string]string{
			"error": err.Error(),
		})
		return
	}

	ctx := r.Context()

	if err := api.DB.CreateTodo(ctx, &todo); err != nil {
		logger.WithError(err).Warn("Error creating todo.")
		utils.WriteError(w, http.StatusInternalServerError, "Error creating todo.", nil)
		return
	}

	logger.WithField("todoID", todo.ID).Info("Todo Created")

	createdTodo, err := api.DB.GetTodoByID(ctx, todo.ID)
	if err != nil {
		logger.WithError(err).Warn("Error getting todo")
		utils.WriteError(w, http.StatusConflict, "Error getting todo", nil)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, createdTodo)
}

type ListResponse struct {
	Todos []*model.Todo `json:"todos"`
}

//GET - /users/{userID}/todos
func (api *TodoAPI) List(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "todo.go -> List()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["userID"])

	logger = logger.WithFields(logrus.Fields{
		"userID":    userID,
	})

	ctx := r.Context()
	todos, err := api.DB.ListTodosByUserID(ctx, userID)
	if err != nil {
		logger.WithError(err).Warn("Error getting todos")
		utils.WriteError(w, http.StatusConflict, "Error getting todos", nil)
		return
	}

	if todos == nil {
		todos = make([]*model.Todo, 0)
	}

	logger.WithField("todosAmount", len(todos)).Info("Todos returned")

	utils.WriteJSON(w, http.StatusOK, &ListResponse{
		Todos: todos,
	})
}

//GET - /users/{userID}/todos/{todoID}
func (api *TodoAPI) Get(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "todo.go -> Get()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["userID"])
	todoID := model.TodoID(vars["todoID"])

	logger = logger.WithFields(logrus.Fields{
		"userID":    userID,
		"todoID": todoID,
	})

	ctx := r.Context()
	todo, err := api.DB.GetTodoByID(ctx, todoID)
	if err != nil {
		logger.WithError(err).Warn("Error getting todo")
		utils.WriteError(w, http.StatusConflict, "Error getting todo", nil)
		return
	}

	logger.Info("Todo returned")

	utils.WriteJSON(w, http.StatusOK, todo)
}

//PATCH -  /users/{userID}/todos/{todoID}
func (api *TodoAPI) Update(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "todo.go -> Update()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["userID"])
	todoID := model.TodoID(vars["todoID"])

	logger = logger.WithFields(logrus.Fields{
		"userID":    userID,
		"todoID": todoID,
	})

	//Decode parameters
	var todoRequest model.Todo
	if err := json.NewDecoder(r.Body).Decode(&todoRequest); err != nil {
		logger.WithError(err).Warn("could not decode parameters")
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	ctx := r.Context()
	todo, err := api.DB.GetTodoByID(ctx, todoID)
	if err != nil {
		logger.WithError(err).Warn("Error getting todo")
		utils.WriteError(w, http.StatusConflict, "Error getting todo", nil)
		return
	}

	if todoRequest.Title != nil ||  len(*todoRequest.Title) != 0 {
		todo.Title =  todoRequest.Title
	}

	if todoRequest.Color != nil ||  len(*todoRequest.Color) != 0 {
		todo.Color =  todoRequest.Color
	}

	if todoRequest.Description != nil ||  len(*todoRequest.Description) != 0 {
		todo.Description =  todoRequest.Description
	}

	if todoRequest.IsFinished != nil {
		todo.IsFinished =  todoRequest.IsFinished
	}

	if err := api.DB.UpdateTodo(ctx, todo); err != nil {
		logger.WithError(err).Warn("Error updating todo.")
		utils.WriteError(w, http.StatusInternalServerError, "Error updating todo.", nil)
		return
	}

	updatedTodo, err := api.DB.GetTodoByID(ctx, todoID)
	if err != nil {
		logger.WithError(err).Warn("Error getting todo")
		utils.WriteError(w, http.StatusConflict, "Error getting todo", nil)
		return
	}

	logger.Info("Todo Updated")

	utils.WriteJSON(w, http.StatusOK, updatedTodo)
}

//DELETE - /users/{userID}/todos/{todoID}
func (api *TodoAPI) Delete(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "todo.go -> Delete()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["userID"])
	todoID := model.TodoID(vars["todoID"])

	logger = logger.WithFields(logrus.Fields{
		"userID":    userID,
		"todoID": todoID,
	})

	ctx := r.Context()
	deleted, err := api.DB.DeleteTodo(ctx, todoID)
	if err != nil {
		logger.WithError(err).Warn("Error deleting todo")
		utils.WriteError(w, http.StatusConflict, "Error deleting todo", nil)
		return
	}

	logger.Info("Todo deleted")

	utils.WriteJSON(w, http.StatusOK, &ActDeleted{
		Deleted: deleted,
	})
}
