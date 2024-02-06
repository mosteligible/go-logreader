package app

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
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
}

func (app *App) status(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJson(w, http.StatusOK, map[string]int16{"status": http.StatusOK})
}
