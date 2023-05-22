package api

import (
	"encoding/json"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type apiFunc func(http.ResponseWriter, *http.Request) error

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        // r.Header.Set("Content-Type", "application/json")
        r.Header.Set("Access-Control-Allow-Origin", "*")

		start := time.Now()

		err := f(w, r)

		if err != nil {
			log.WithFields(log.Fields{
				"method":   r.Method,
				"path":     r.URL.Path,
				"addr":     r.RemoteAddr,
				"duration": time.Since(start),
			}).Errorf("%s", err.Error())

			WriteJSON(w, http.StatusBadRequest, err.Error())
			return
		}

		log.WithFields(log.Fields{
			"method":   r.Method,
			"path":     r.URL.Path,
			"addr":     r.RemoteAddr,
			"duration": time.Since(start),
		}).Infof("%s %s", r.Method, r.URL.Path)
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
