package handler

import (
	"encoding/json"
	"github.com/DCloudGaming/cloud-morph-host/pkg/env"
	"github.com/DCloudGaming/cloud-morph-host/pkg/jwt"
	"github.com/DCloudGaming/cloud-morph-host/pkg/model"
	"github.com/DCloudGaming/cloud-morph-host/pkg/utils"
	"net/http"
	"os"
)

func ApiHandlerWrapper(
	sharedEnv *env.SharedEnv,
	f_in func(
		sharedEnv *env.SharedEnv, w http.ResponseWriter, r *http.Request,
		u *model.User, hostU *model.User, headPath string)) (f func(w http.ResponseWriter, r *http.Request)) {
	return func(w http.ResponseWriter, r *http.Request) {

		// Reply preflight query to avoid CORS blocking
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("WEB_PROTOCOL") + "://" + os.Getenv("WEB_HOST"))
		//w.Header().Set("Access-Control-Allow-Origin", "https://www.declo.co")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Content-Length, Authorization, Accept,X-Requested-With,Origin")
		if r.Method == http.MethodOptions {
			json.NewEncoder(w).Encode("OKOK")
		}

		u, _ := jwt.HandleUserCookie(*sharedEnv, w, r)

		//decoder := json.NewDecoder(r.Body)
		//var req model.HostJwtToken
		//err := decoder.Decode(&req)
		var hostU *model.User = nil
		//if err == nil && &req != nil {
		//	hostU, _ = jwt.DecodeUser(req.Token)
		//}

		var head string
		head, r.URL.Path = utils.ShiftPath(r.URL.Path)
		head, r.URL.Path = utils.ShiftPath(r.URL.Path)
		head, r.URL.Path = utils.ShiftPath(r.URL.Path)

		f_in(sharedEnv, w, r, u, hostU, head)
	}
}
