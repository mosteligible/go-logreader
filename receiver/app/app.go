package app

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mosteligible/go-logreader/receiver/config"
	"github.com/mosteligible/go-logreader/receiver/core/broker"
	"github.com/mosteligible/go-logreader/receiver/core/logstream"
	"github.com/mosteligible/go-logreader/receiver/core/middlewares"
	"github.com/mosteligible/go-logreader/receiver/core/utils"
)

type App struct {
	Router *mux.Router
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

	subRouter := app.Router.PathPrefix("/log").Subrouter()
	subRouter.HandleFunc("/logit", app.logStream).Methods(http.MethodPost)

	subRouter.Use(middlewares.ApiKey)
}

func (app *App) Run() {
	log.Printf("Listening in port: %s\n", config.Env.AppPort)
	log.Fatal(http.ListenAndServe(config.Env.AppPort, app.Router))
}

func (app *App) status(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJson(w, http.StatusOK, map[string]int16{"status": http.StatusOK})
}

func (app *App) logStream(w http.ResponseWriter, r *http.Request) {
	log.Println("POST - /logit")
	var logMsg logstream.LogStream
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&logMsg); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error(), true)
		return
	}
	conn := broker.GetConnection(logMsg.ClientId)
	if err := utils.SendMsgWithRetries(logMsg.Message, conn); err != nil {
		utils.RespondWithError(
			w, http.StatusInternalServerError, "Connection issues with broker", true,
		)
		return
	}

	utils.RespondWithJson(
		w, http.StatusAccepted, map[string]int16{"status": http.StatusAccepted},
	)
}
