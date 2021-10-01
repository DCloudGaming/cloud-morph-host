package handler

import (
	"encoding/json"
	"github.com/DCloudGaming/cloud-morph-host/pkg/env"
	"github.com/DCloudGaming/cloud-morph-host/pkg/jwt"
	"github.com/DCloudGaming/cloud-morph-host/pkg/model"
	"github.com/DCloudGaming/cloud-morph-host/pkg/utils"
	"net/http"
)

func ApiHandlerWrapper(
	sharedEnv *env.SharedEnv,
	f_in func(
		sharedEnv *env.SharedEnv, w http.ResponseWriter, r *http.Request,
		u *model.User, headPath string)) (f func(w http.ResponseWriter, r *http.Request)) {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Credentials", "true")
		// Reply preflight query to avoid CORS blocking

		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Content-Length, Authorization, Accept,X-Requested-With,Origin")
			//w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if r.Method == http.MethodOptions {
			json.NewEncoder(w).Encode("OKOK")
			//return
		}

		u, _ := jwt.HandleUserCookie(*sharedEnv, w, r)

		var head string
		head, r.URL.Path = utils.ShiftPath(r.URL.Path)
		head, r.URL.Path = utils.ShiftPath(r.URL.Path)
		head, r.URL.Path = utils.ShiftPath(r.URL.Path)

		f_in(sharedEnv, w, r, u, head)
	}
}
