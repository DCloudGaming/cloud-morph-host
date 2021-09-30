package handler

import (
	"encoding/json"
	"github.com/DCloudGaming/cloud-morph-host/pkg/env"
	"github.com/DCloudGaming/cloud-morph-host/pkg/errors"
	"github.com/DCloudGaming/cloud-morph-host/pkg/model"
	"github.com/DCloudGaming/cloud-morph-host/pkg/perm"
	"github.com/DCloudGaming/cloud-morph-host/pkg/write"
	"net/http"
	"strconv"
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
		} else if head == "discover" {
			getDiscoverApps(*sharedEnv, *u, w, r)
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

	isAllow := perm.RequireAuthenticated(sharedEnv, w, r)
	if !isAllow {
		write.Error(errors.RouteUnauthorized, w, r)
		return
	}

	dbHostApps, _ := sharedEnv.AppRepo().GetFromHost(walletAddress)
	write.JSON(dbHostApps, w, r)
}

func getDiscoverApps(sharedEnv env.SharedEnv, u model.User, w http.ResponseWriter, r *http.Request) {
	dbHostApps, _ := sharedEnv.AppRepo().GetAllRegisteredApps()
	write.JSON(dbHostApps, w, r)
}

func registerApp(sharedEnv env.SharedEnv, u model.User, w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req model.RegisterAppReq
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

	rowsAffected, _ := sharedEnv.AppRepo().RegisterBatch(req)
	write.JSON(model.RegisterBatchResponse{RowsAffected: rowsAffected}, w, r)
}

func getSession(sharedEnv env.SharedEnv, u model.User, w http.ResponseWriter, r *http.Request) {
	sessionId, _ := strconv.Atoi(r.URL.Query().Get("session_id"))

	isAllow := perm.RequireAuthenticated(sharedEnv, w, r)
	if !isAllow {
		write.Error(errors.RouteUnauthorized, w, r)
		return
	}

	session, _ := sharedEnv.StreamSessionRepo().GetSession(sessionId)
	write.JSON(session, w, r)
}

// For now , let only client side initiate the session, the player side will pull this status to
// verify session really starts
func startSession(sharedEnv env.SharedEnv, u model.User, w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req model.StartSessionReq
	err := decoder.Decode(&req)
	if err != nil || &req == nil {
		write.Error(errors.NoJSONBody, w, r)
		return
	}

	isAllow := perm.RequireOwner(u.WalletAddress, req.ClientWalletAddress) &&
		perm.RequireAuthenticated(sharedEnv, w, r)
	if !isAllow {
		write.Error(errors.RouteUnauthorized, w, r)
		return
	}

	session, _ := sharedEnv.StreamSessionRepo().StartSession(req)
	write.JSON(session, w, r)
}

func updateSession(sharedEnv env.SharedEnv, u model.User, w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req model.UpdateSessionReq
	err := decoder.Decode(&req)
	if err != nil || &req == nil {
		write.Error(errors.NoJSONBody, w, r)
		return
	}

	dbSession, _ := sharedEnv.StreamSessionRepo().GetSession(req.SessionID)

	isAllow := perm.RequireOwner(u.WalletAddress, dbSession.ClientWalletAddress) &&
		perm.RequireAuthenticated(sharedEnv, w, r)
	if !isAllow {
		write.Error(errors.RouteUnauthorized, w, r)
		return
	}

	sharedEnv.StreamSessionRepo().UpdateSession(req)
}



