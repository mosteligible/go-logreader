package app

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mosteligible/go-logreader/box/config"
	"github.com/mosteligible/go-logreader/box/core/utils"
)

type App struct {
	Router *mux.Router
}

func NewApp() App {
	router := mux.NewRouter()
	app := App{Router: router}
	app.initialize()
	return app
}

func (app *App) initialize() {
	app.Router.HandleFunc("/status", app.status).Methods(http.MethodGet)
}

func (app *App) Run() {
	log.Printf("Listening to port: %s", config.Env.AppPort)
	log.Fatal(http.ListenAndServe(config.Env.AppPort, app.Router))
}

func (app *App) status(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJson(w, http.StatusOK, map[string]int16{"status": http.StatusOK})
}
