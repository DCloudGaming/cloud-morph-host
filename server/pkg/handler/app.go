package handler

import (
	"encoding/json"
	"github.com/DCloudGaming/cloud-morph-host/pkg/env"
	"github.com/DCloudGaming/cloud-morph-host/pkg/errors"
	"github.com/DCloudGaming/cloud-morph-host/pkg/model"
	"github.com/DCloudGaming/cloud-morph-host/pkg/perm"
	"github.com/DCloudGaming/cloud-morph-host/pkg/schema"
	"github.com/DCloudGaming/cloud-morph-host/pkg/write"
	"net/http"
)

func AppHandler(
	sharedEnv *env.SharedEnv, w http.ResponseWriter, r *http.Request,
	u *model.User, head string) {
	switch r.Method {
	case http.MethodGet:
		if head == "" {
			getHostApps(*sharedEnv, *u, w, r)
		} else if head == "getSession" {
			getSession(*sharedEnv, *u, w, r)
		} else if head == "getSessions" {
			getSessions(*sharedEnv, *u, w, r)
		} else {
			write.Error(errors.RouteNotFound, w, r)
		}
	case http.MethodPost:
		if head == "registerApp" {
			registerApp(*sharedEnv, *u, w, r)
		} else if head == "startSession" {
			startSession(*sharedEnv, *u, w, r)
		} else if head == "updateSession" {
			updateSession(*sharedEnv, *u, w, r)
		} else {
			write.Error(errors.RouteNotFound, w, r)
		}
	default:
		write.Error(errors.BadRequestMethod, w, r)
	}
}

func getHostApps(sharedEnv env.SharedEnv, u model.User, w http.ResponseWriter, r *http.Request) {
	walletAddress := r.URL.Query().Get("wallet_address")

	dbHostApps, _ := sharedEnv.AppRepo().GetFromHost(walletAddress)
	write.JSON(dbHostApps, w, r)
}

func registerApp(sharedEnv env.SharedEnv, u model.User, w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req schema.RegisterAppReq
	err := decoder.Decode(&req)
	if err != nil || &req == nil {
		write.Error(errors.NoJSONBody, w, r)
		return
	}

	isAllow := perm.RequireOwner(u.WalletAddress, req.WalletAddress) &&
		perm.RequireAuthenticated(sharedEnv, w, r)
	if !isAllow {
		write.Error(errors.RouteUnauthorized, w, r)
		return
	}

	dbRegisterApps, _ := sharedEnv.AppRepo().RegisterBatch(req)
	write.JSON(dbRegisterApps, w, r)
}

// For now , let only client side initiate the session, the player side will pull this status to
// verify session really starts
func startSession(sharedEnv env.SharedEnv, u model.User, w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req schema.StartSessionReq
	err := decoder.Decode(&req)
	if err != nil || &req == nil {
		write.Error(errors.NoJSONBody, w, r)
		return
	}

	isAllow := perm.RequireOwner(u.WalletAddress, req.HostWalletAddress) &&
		perm.RequireAuthenticated(sharedEnv, w, r)
	if !isAllow {
		write.Error(errors.RouteUnauthorized, w, r)
		return
	}

	dbRegisterApps, _ := sharedEnv.AppRepo().StartSession(req)
	write.JSON(dbRegisterApps, w, r)
}

func updateSession(sharedEnv env.SharedEnv, u model.User, w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req schema.UpdateSessionReq
	err := decoder.Decode(&req)
	if err != nil || &req == nil {
		write.Error(errors.NoJSONBody, w, r)
		return
	}

	dbSession, _ := sharedEnv.AppRepo().GetSession(req.SessionID)

	isAllow := perm.RequireOwner(u.WalletAddress, dbSession.HostWalletAddress) &&
		perm.RequireAuthenticated(sharedEnv, w, r)
	if !isAllow {
		write.Error(errors.RouteUnauthorized, w, r)
		return
	}


	dbRegisterApps, _ := sharedEnv.AppRepo().UpdateSession(req)
	write.JSON(dbRegisterApps, w, r)
}




