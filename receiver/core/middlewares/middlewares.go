package middlewares

import (
	"log"
	"net/http"

	"github.com/mosteligible/go-logreader/receiver/config"
)

func ApiKey(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api_key := r.Header.Get("api-key")
		if api_key != config.CLIENT_DATA_API_KEY {
			log.Printf("forbidden header api-key: %s\n", r.Header.Get("api-key"))

			http.Error(w, "Forbidded api-key", http.StatusForbidden)
			return
		}
		h.ServeHTTP(w, r)
	})
}
