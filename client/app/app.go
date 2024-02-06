package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
	"github.com/mosteligible/go-logreader/client/config"
	"github.com/mosteligible/go-logreader/client/core/middlewares"
	"github.com/mosteligible/go-logreader/client/customer"
)

type App struct {
	Router       *mux.Router
	statusRouter *mux.Router
	DB           *sql.DB
}

func NewApp() App {
	app := App{}
	app.initialize()
	return app
}

func (app *App) initialize() {
	if !config.SslOk {
		config.DB_SSL_MODE = "disable"
	}
	db_conn_string := fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=%s host=%s",
		config.CLIENT_DB_USERNAME,
		config.CLIENT_DB_PASSWORD,
		config.POSTGRES_DB,
		config.DB_SSL_MODE,
		config.CLIENT_DB_HOST,
	)
	var err error
	app.DB, err = sql.Open("postgres", db_conn_string)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to database..")

	app.Router = mux.NewRouter().StrictSlash(false)
	app.initializeRoutes()
}

func (app *App) initializeRoutes() {
	app.Router.HandleFunc("/status", app.status).Methods("GET")
	prefixedRouter := app.Router.PathPrefix("/customers").Subrouter()
	prefixedRouter.HandleFunc("", app.getCustomers).Methods("GET")
	prefixedRouter.HandleFunc("", app.addCustomer).Methods("POST")
	prefixedRouter.HandleFunc("", app.updateCustomer).Methods("PUT")
	prefixedRouter.HandleFunc("/{id:[a-zA-Z]+}", app.getCustomer).Methods("GET")
	prefixedRouter.HandleFunc("/{id:[a-zA-Z]+}", app.deleteCustomer).Methods("DELETE")
	prefixedRouter.Use(middlewares.ApiKey)
}

func (app *App) Run() {
	log.Printf("Listening in port: %s\n", config.APP_PORT)
	log.Fatal(http.ListenAndServe(config.APP_PORT, app.Router))
}

func (app *App) status(w http.ResponseWriter, r *http.Request) {
	log.Println("GET - /status - 200")
	respondWithJson(w, http.StatusOK, map[string]int{"status": 200})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJson(w, code, map[string]string{"error": message})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (app *App) addCustomer(w http.ResponseWriter, r *http.Request) {
	var customer customer.Customer
	log.Println("POST - /customers ")
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&customer); err != nil {
		log.Printf("Customer: %v", customer)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := customer.AddCustomer(app.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// update services with new customer data
	// go notifier.NotifyService()
	respondWithJson(w, http.StatusCreated, customer)
}

func (app *App) deleteCustomer(w http.ResponseWriter, r *http.Request) {
	var cust customer.Customer
	log.Println("DELETE - /customers")
	vars := mux.Vars(r)
	custId := vars["id"]
	cust = customer.Customer{Id: custId}
	if err := cust.DeleteCustomer(app.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// update services with new customer data
	// go notifier.NotifyService()
	respondWithJson(w, http.StatusAccepted, map[string]int{"status": http.StatusAccepted})
}

func (app *App) updateCustomer(w http.ResponseWriter, r *http.Request) {
	var cust customer.Customer
	log.Println("PUT - /customers")
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&cust); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := cust.UpdateCustomer(app.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusAccepted, cust)
}

func (app *App) getCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	custId := vars["id"]
	log.Printf("GET - /customers/%s", custId)
	customer, err := customer.GetCustomer(app.DB, custId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, customer)
}

func (app *App) getCustomers(w http.ResponseWriter, r *http.Request) {
	log.Println("GET - /customers")
	customers, err := customer.GetAllCustomers(app.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, customers)
}
