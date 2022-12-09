package main

import (
	"net/http"
	cTodo "wb/app/controller/todo"
	hTodo "wb/app/handler/todo"

	cAuth "wb/app/controller/auth"
	hAuth "wb/app/handler/auth"

	"wb/app/model"
	"wb/db"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

func main() {
	dbConn, err := db.ConnectPG()
	if err != nil {
		panic(err)
	}

	err = dbConn.AutoMigrate(&model.Todo{}, &model.Auth{})
	if err != nil {
		panic(err)
	}

	redisClient = db.ConnectRedis()

	todoController := cTodo.NewTodo(dbConn)
	todoHandler := hTodo.Handler{
		Contract: todoController,
	}

	authController := cAuth.NewAuth(redisClient, dbConn)
	authHandler := hAuth.Handler{
		Contract: authController,
	}

	http.Handle("/todo", Auth(RequestMethodGet(http.HandlerFunc(todoHandler.Todo)), false))
	http.Handle("/todo/get", RequestMethodGet(http.HandlerFunc(todoHandler.Get)))
	http.Handle("/todo/insert", RequestMethodPost(http.HandlerFunc(todoHandler.Insert)))
	http.Handle("/todo/update", RequestMethodPost(http.HandlerFunc(todoHandler.Update)))
	http.Handle("/todo/delete", RequestMethodPost(http.HandlerFunc(todoHandler.Delete)))

	http.Handle("/auth/register", RequestMethodGet(http.HandlerFunc(authHandler.RegisterPage)))
	http.Handle("/auth/login", Auth(RequestMethodGet((http.HandlerFunc(authHandler.LoginPage))), true))
	http.Handle("/auth/doRegister", RequestMethodPost(http.HandlerFunc(authHandler.DoRegister)))
	http.Handle("/auth/doLogin", RequestMethodPost(http.HandlerFunc(authHandler.DoLogin)))
	http.Handle("/auth/logout", RequestMethodGet(http.HandlerFunc(authHandler.Logout)))

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

func Auth(next http.Handler, isLoginPage bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("SESSION_TOKEN")
		if err != nil {
			if isLoginPage {
				next.ServeHTTP(w, r)
			} else {
				http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
			}
		} else {
			if isLoginPage {
				http.Redirect(w, r, "/todo", http.StatusSeeOther)
			} else {
				next.ServeHTTP(w, r)
			}
		}

	})
}
