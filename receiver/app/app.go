package app

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mosteligible/go-logreader/receiver/config"
	"github.com/mosteligible/go-logreader/receiver/core/logstream"
	"github.com/mosteligible/go-logreader/receiver/core/middlewares"
	"github.com/mosteligible/go-logreader/receiver/core/utils"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func NewApp() App {
	app := App{}
	app.initialize()
	return app
}

func (app *App) initialize() {
	app.Router = mux.NewRouter()
	app.initializeRoutes()
}

func (app *App) initializeRoutes() {
	app.Router.HandleFunc("/status", app.status).Methods(http.MethodGet)

	subRouter := app.Router.NewRoute().Subrouter()

	subRouter.Use(middlewares.ApiKey)
}

func (app *App) Run() {
	log.Printf("Listening in port: %s\n", config.APP_PORT)
	log.Fatal(http.ListenAndServe(config.APP_PORT, app.Router))
}

func (app *App) status(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJson(w, http.StatusOK, map[string]int16{"status": http.StatusOK})
}

func (app *App) logStream(w http.ResponseWriter, r *http.Request) {
	log.Println("POST - /message")
	var logMsg logstream.LogStream
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&logMsg); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error(), true)
		return
	}
	utils.RespondWithJson(
		w, http.StatusAccepted, map[string]int16{"status": http.StatusAccepted},
	)
}
