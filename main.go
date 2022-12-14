package main

import (
	"context"
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
	http.Handle("/todo/get", Auth(RequestMethodGet(http.HandlerFunc(todoHandler.Get)), false))
	http.Handle("/todo/insert", Auth(RequestMethodPost(http.HandlerFunc(todoHandler.Insert)), false))
	http.Handle("/todo/update", Auth(RequestMethodPost(http.HandlerFunc(todoHandler.Update)), false))
	http.Handle("/todo/delete", Auth(RequestMethodPost(http.HandlerFunc(todoHandler.Delete)), false))

	http.Handle("/auth/register", Auth(RequestMethodGet(http.HandlerFunc(authHandler.RegisterPage)), true))
	http.Handle("/auth/login", Auth(RequestMethodGet((http.HandlerFunc(authHandler.LoginPage))), true))
	http.Handle("/auth/doRegister", Auth(RequestMethodPost(http.HandlerFunc(authHandler.DoRegister)), true))
	http.Handle("/auth/doLogin", Auth(RequestMethodPost(http.HandlerFunc(authHandler.DoLogin)), true))
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

func Auth(next http.Handler, isAuthPage bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("SESSION_TOKEN")
		if err != nil {
			if isAuthPage {
				next.ServeHTTP(w, r)
				return
			}

			http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
			return
		}

		if isAuthPage {
			http.Redirect(w, r, "/todo", http.StatusSeeOther)
			return
		}

		userID, err := redisClient.Get(context.TODO(), cookie.Value).Result()
		if err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:   "SESSION_TOKEN",
				Path:   "/",
				Value:  "",
				MaxAge: -1,
			})
			http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
			return
		}
		ctx := context.WithValue(r.Context(), "user_id", userID)

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
