package handler

import (
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

		u, _ := jwt.HandleUserCookie(*sharedEnv, w, r)

		var head string
		head, r.URL.Path = utils.ShiftPath(r.URL.Path)
		head, r.URL.Path = utils.ShiftPath(r.URL.Path)
		head, r.URL.Path = utils.ShiftPath(r.URL.Path)

		f_in(sharedEnv, w, r, u, head)
	}
}
