package write

import (
	"encoding/json"
	"log"
	"net/http"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/DCloudGaming/cloud-morph-host/pkg/errors"
)

type errorResponse struct {
	Error string
}

func Error(err error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}

func JSON(obj interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(obj)
	}
}

func JSONorErr(obj interface{}, err error) http.HandlerFunc {
	if err != nil {
		return Error(err)
	}

	return JSON(obj)
}

func StreamFile(file_name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pwd, _ := os.Getwd()
		fmt.Println("WORKING DIR")
		fmt.Println(pwd)
		full_file_name := pwd + "/" + file_name

		streamFilebytes, err := ioutil.ReadFile(full_file_name)

		if err != nil {
			fmt.Println(err)
		}

		b := bytes.NewBuffer(streamFilebytes)

		w.Header().Set("Content-type", "application/octet-stream")

		if _, err := b.WriteTo(w); err != nil { // <----- here!
			fmt.Fprintf(w, "%s", err)
		}

		w.Write([]byte("File returned"))
	}
}
