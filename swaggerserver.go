package main

import (
	"net/http"

	_ "wb/docs"

	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"
)

func ServeSwagger() {
	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:1323/swagger/doc.json"), //The url pointing to API definition
	))
	http.ListenAndServe(":1323", r)
}
