package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/suarezgary/GolangApi/config"
	v1 "github.com/suarezgary/GolangApi/controllers/v1"
	"github.com/suarezgary/GolangApi/middleware"
	"github.com/suarezgary/GolangApi/utils/httpmiddleware"
)

// Log Log
var Log = config.Cfg().GetLogger()

// CreateRouter CreateRouter
func CreateRouter() http.Handler {
	router := mux.NewRouter()
	router.StrictSlash(true)

	// V1 Routes
	v1Router := router.PathPrefix("/v1").Subrouter()
	v1Router.HandleFunc("/", v1.API).Methods("GET")
	v1Router.HandleFunc("/gophers", httpmiddleware.Use(v1.GetGophers, middleware.RequireAPIKey)).Methods("GET")

	return httpmiddleware.Use(router.ServeHTTP, middleware.GetContext, httpmiddleware.RecoverInternalServerError)
}
