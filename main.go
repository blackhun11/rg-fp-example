package main

import (
	"context"
	"log"
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

//	@title			Swagger GENERATE KE 3
//	@version		1.0
//	@description	TODO Server
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Ferianto Surya Wijaya
//	@contact.url	https://feriantosw.my.id
//	@contact.email	feriantosw77@gmail.com

//	@license.name				Apache 2.0
//	@license.url				http://www.apache.org/licenses/LICENSE-2.0.html
//	@securityDefinitions.basic	BasicAuth
//	@in							header
//	@name						Authorization
//	@host						localhost:8080
//	@BasePath					/
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

	// mencegah untuk yang belum login ke todo
	http.Handle("/todo", AllowCORS(Auth(RequestMethodGet(http.HandlerFunc(todoHandler.Todo)), false)))
	http.Handle("/todo/get", AllowCORS(Auth(RequestMethodGet(http.HandlerFunc(todoHandler.Get)), false)))
	http.Handle("/todo/insert", AllowCORS(Auth(RequestMethodPost(http.HandlerFunc(todoHandler.Insert)), false)))
	http.Handle("/todo/update", AllowCORS(Auth(RequestMethodPost(http.HandlerFunc(todoHandler.Update)), false)))
	http.Handle("/todo/delete", AllowCORS(Auth(RequestMethodPost(http.HandlerFunc(todoHandler.Delete)), false)))

	// mencegah untuk yang udah login ke auth
	http.Handle("/auth/register", AllowCORS(Auth(RequestMethodGet(http.HandlerFunc(authHandler.RegisterPage)), true)))
	http.Handle("/auth/login", AllowCORS(Auth(RequestMethodGet((http.HandlerFunc(authHandler.LoginPage))), true)))
	http.Handle("/auth/doRegister", AllowCORS(Auth(RequestMethodPost(http.HandlerFunc(authHandler.DoRegister)), true)))
	http.Handle("/auth/doLogin", AllowCORS(Auth(RequestMethodPost(http.HandlerFunc(authHandler.DoLogin)), true)))
	http.Handle("/auth/logout", AllowCORS(RequestMethodGet(http.HandlerFunc(authHandler.Logout))))

	http.Handle("/healthcheck", AllowCORS(http.HandlerFunc(healthcheck)))
	go ServeSwagger()

	log.Println("Application Running")
	http.ListenAndServe(":8080", nil)
}

//	@Summary		Healthcheck
//	@Description	Do Healthcheck
//	@Tags			Healthcheck
//	@Accept			json
//	@Produce		plain
//	@Success		200	string	healthy
//	@Router			/healthcheck [get]
func healthcheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("healthy"))
}

func AllowCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
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
		// no cookie --> redirect to login
		if err != nil {
			if isAuthPage {
				next.ServeHTTP(w, r)
				return
			}
			// redirect to login
			http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
			return
		}

		// have cookie but auth page --> redirect to todo
		if isAuthPage {
			http.Redirect(w, r, "/todo", http.StatusSeeOther)
			return
		}

		// cookie:user_id
		// have cookie --> do cookie validation
		userID, err := redisClient.Get(context.TODO(), cookie.Value).Result()

		// have error --> cookie not valid, delete the cookie
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

		// cookie valid --> pass the user id via context
		ctx := context.WithValue(r.Context(), "user_id", userID)

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
