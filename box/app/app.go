package app

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mosteligible/go-logreader/box/config"
	"github.com/mosteligible/go-logreader/box/constants"
	"github.com/mosteligible/go-logreader/box/core/broker"
	"github.com/mosteligible/go-logreader/box/core/models"
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
	app.Router.HandleFunc("/clientpool", app.clientUpdates).Methods(http.MethodPost)
}

func (app *App) Run() {
	log.Printf("Listening to port: %s", config.Env.AppPort)
	log.Fatal(http.ListenAndServe(config.Env.AppPort, app.Router))
}

func (app *App) status(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJson(w, http.StatusOK, map[string]int16{"status": http.StatusOK})
}

func (app *App) clientUpdates(w http.ResponseWriter, r *http.Request) {
	var client models.Customer
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&client); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error(), true)
		return
	}

	event := client.Event
	if event == "" {
		utils.RespondWithError(
			w, http.StatusBadRequest, "no client update event type defined", true,
		)
		return
	}

	switch event {
	case constants.ClientAdded:
		conn := broker.NewConnection(client.Id, client.Id, []string{})
		go conn.Consume()
	case constants.ClientRemoved:
		if _, ok := broker.ConnectionPool[client.Id]; !ok {
			utils.RespondWithError(
				w, http.StatusNotFound, "Not available", true,
			)
			return
		}
		delete(broker.ConnectionPool, client.Id)
	case constants.ClientUpdated:
		conn := broker.NewConnection(client.Id, client.Id, []string{})
		go conn.Consume()
	default:
		utils.RespondWithError(
			w, http.StatusNotFound, "invalid update event type", true,
		)
		return
	}
	utils.RespondWithJson(w, http.StatusAccepted, map[string]string{"test": "ok"})
}
