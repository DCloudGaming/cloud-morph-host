package perm

import (
	"github.com/DCloudGaming/cloud-morph-host/pkg/env"
	"github.com/DCloudGaming/cloud-morph-host/pkg/jwt"
	"github.com/DCloudGaming/cloud-morph-host/pkg/model"
	"net/http"
)

const (
	RegisterAppType int = 0
	StartSessionType = 1
	GetHostAppsType = 2
)

// TODO: find ways to cast dynamically the type of req schema
//func extractReqSchema(w http.ResponseWriter, r *http.Request, reqType int) (interface{}, bool) {
//	decoder := json.NewDecoder(r.Body)
//	switch reqType {
//	case RegisterAppType:
//		var req schema.RegisterAppReq
//	default:
//		var req interface{}
//	}
//	err := decoder.Decode(&req)
//}
//

func RequireAuthenticated(sharedEnv env.SharedEnv, w http.ResponseWriter, r *http.Request) (bool) {
	_, allowJwt := jwt.RequireAuth(model.StatusUnverified, sharedEnv, w, r)
	return allowJwt
}

func RequireOwner(address1 string,address2 string) (bool) {
	if address1 != address2 {
		return false
	}
	return true
}

func RequireAdmin(sharedEnv env.SharedEnv, checkAddress string) (bool) {
	isAdmin := sharedEnv.UserRepo().VerifyAdmin(checkAddress)
	return isAdmin
}