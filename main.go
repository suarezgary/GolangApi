package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/suarezgary/GolangApi/config"
	"github.com/suarezgary/GolangApi/routes"
)

//Log logger
var Log = config.Cfg().GetLogger()

// CORSHandler This CORS Handler comes with some pretty lenient defaults, depending on your application,
// you may want to curtail some of these open settings
var CORSHandler = handlers.CORS(
	handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
	handlers.AllowCredentials(), handlers.AllowedHeaders(
		[]string{"x-locale", "x-api-key", "content-type",
			"access-control-request-headers", "access-control-request-method",
			"x-csrftoken"}),
	handlers.AllowedOrigins(config.Cfg().AllowedOrigins))

func main() {
	fmt.Println("test")
	Log.Info("Setting up database connection ...")

	/*for {
		err := models.Setup()
		if err != nil {
			Log.WithError(err).Error("Error setting up database connection, retrying ...")
			time.Sleep(time.Second * 3)
		} else {
			break
		}
	}*/

	Log.Info("Connected to database")

	// TODO: Implement graceful stop
	Log.Info("Starting HTTP server")
	http.ListenAndServe(
		fmt.Sprintf("%s:%s", config.Cfg().ListenAddress, config.Cfg().ListenPort),
		CORSHandler(handlers.CombinedLoggingHandler(os.Stdout, routes.CreateRouter())),
	)
	Log.Info("Stopped HTTP server")
}
