package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mosteligible/go-logreader/box/core/models"
	"github.com/mosteligible/go-logreader/box/core/utils"
)

func CustomerUpdates(w http.ResponseWriter, r *http.Request) {
	var customer models.Customer
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&customer); err != nil {
		log.Printf("Error <%s> while decoding post body.", err.Error())
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error(), true)
		return
	}
	defer r.Body.Close()
}
