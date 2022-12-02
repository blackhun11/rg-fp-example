package main

import (
	"net/http"
	cTodo "wb/app/controller/todo"
	"wb/app/handler/todo"

	"wb/db"
)

func main() {
	err := db.Connect()
	if err != nil {
		panic(err)
	}

	err = db.Migrate()
	if err != nil {
		panic(err)
	}
	todoController := cTodo.NewController(db.DBConn)
	todoHandler := todo.Handler{
		Contract: todoController,
	}
	http.Handle("/todo", RequestMethodGet(http.HandlerFunc(todoHandler.Todo)))
	http.Handle("/todo/get", RequestMethodGet(http.HandlerFunc(todoHandler.Get)))
	http.Handle("/todo/insert", RequestMethodPost(http.HandlerFunc(todoHandler.Insert)))
	http.Handle("/todo/update", RequestMethodPost(http.HandlerFunc(todoHandler.Update)))
	http.Handle("/todo/delete", RequestMethodPost(http.HandlerFunc(todoHandler.Delete)))

	http.ListenAndServe(":8080", nil)
}

func RequestMethodGet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		if method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method is not allowed"))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func RequestMethodPost(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		if method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method is not allowed"))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
