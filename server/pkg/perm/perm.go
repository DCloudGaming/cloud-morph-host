package perm

import (
	"github.com/DCloudGaming/cloud-morph-host/pkg/env"
	"github.com/DCloudGaming/cloud-morph-host/pkg/jwt"
	"github.com/DCloudGaming/cloud-morph-host/pkg/model"
	"net/http"
	"strings"
)

const (
	RegisterAppType int = 0
	StartSessionType = 1
	GetHostAppsType = 2
)

func RequireAuthenticated(sharedEnv env.SharedEnv, w http.ResponseWriter, r *http.Request) (bool) {
	if (sharedEnv.Mode() == "DEBUG") {
		return true
	}
	_, allowJwt := jwt.RequireAuth(model.StatusUnverified, sharedEnv, w, r)
	return allowJwt
}

func RequireOwner(sharedEnv env.SharedEnv, address1 string, address2 string) (bool) {
	if (sharedEnv.Mode() == "DEBUG") {
		return true
	}
	if strings.ToLower(address1) != strings.ToLower(address2) {
		return false
	}
	return true
}

func RequireAdmin(sharedEnv env.SharedEnv, checkAddress string) (bool) {
	if (sharedEnv.Mode() == "DEBUG") {
		return true
	}
	isAdmin := sharedEnv.UserRepo().VerifyAdmin(strings.ToLower(checkAddress))
	return isAdmin
}