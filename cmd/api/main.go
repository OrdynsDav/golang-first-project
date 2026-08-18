package main

import (
	"log"
	"net/http"
	"os"
	"restapi-tasks/internal/database"
	"restapi-tasks/internal/handlers"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://taskuser:taskpass@localhost:5432/tasksdb?sslmode=disable"
	}

	serverPort := os.Getenv("SERVER_PORT")

	if serverPort == "" {
		serverPort = "8080"
	}

	log.Printf("Starting server on port %s", serverPort)

	db, err := database.Connect(databaseURL)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	log.Println("Success connect to database")

	taskStore := database.NewTaskStore(db)

	handler := handlers.NewHandlerStore(taskStore)

	mux := http.NewServeMux()

	mux.HandleFunc("/tasks", methodHendler(handler.GetAllTasks, "GET"))
	mux.HandleFunc("/tasks/create", methodHendler(handler.CreateTask, "POST"))
	mux.HandleFunc("/tasks/", TaskIDHandler(handler))

	loggedMux := loggingMiddleWare(mux)
	serverAddr := ":" + serverPort

	err = http.ListenAndServe(serverAddr, loggedMux)

	if err != nil {
		log.Fatal(err)
	}
}

func methodHendler(hadnlerFunc http.HandlerFunc, allowedMethod string) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		if request.Method != allowedMethod {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
		hadnlerFunc(w, request)
	}
}

func TaskIDHandler(handler *handlers.Handlers) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			handler.GetTask(w, request)
		case http.MethodPut:
			handler.UpdateTask(w, request)
		case http.MethodDelete:
			handler.DeleteTask(w, request)
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	}
}

func loggingMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		log.Printf("Method: %s,  URLPath: %s, RemoteAddr: %s", request.Method, request.URL.Path, request.RemoteAddr)
		next.ServeHTTP(w, request)
	})
}
