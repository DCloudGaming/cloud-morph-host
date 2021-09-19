package write

import (
	"encoding/json"
	"github.com/DCloudGaming/cloud-morph-host/pkg/errors"
	"log"
	"net/http"
)

type errorResponse struct {
	Error string
}

func Error(err error, w http.ResponseWriter, r *http.Request) {
	found, code := errors.GetCode(err)
	if !found {
		// unexpected error - we should clean this up to avoid showing sql errors in the browser
		log.Println("Unexpected Error: ", err)
		err = errors.InternalError
	}
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&errorResponse{Error: err.Error()})
}

func JSON(obj interface{}, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(obj)
}
