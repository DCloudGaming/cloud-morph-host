package handler

import (
	"encoding/json"
	"github.com/DCloudGaming/cloud-morph-host/pkg/env"
	"github.com/DCloudGaming/cloud-morph-host/pkg/errors"
	"github.com/DCloudGaming/cloud-morph-host/pkg/jwt"
	"github.com/DCloudGaming/cloud-morph-host/pkg/model"
	"github.com/DCloudGaming/cloud-morph-host/pkg/perm"
	"github.com/DCloudGaming/cloud-morph-host/pkg/write"
	"net/http"
)

func UserHandler(
	sharedEnv *env.SharedEnv, w http.ResponseWriter, r *http.Request,
	u *model.User, head string) {
	switch r.Method {
		case http.MethodGet:
			if head == "" {
				getUser(*sharedEnv, w, r)
			} else if head == "profile" {
				getProfile(*sharedEnv, *u, w, r)
			} else {
				write.Error(errors.RouteNotFound, w, r)
			}
		case http.MethodPost:
			if head == "signup" {
				signUp(*sharedEnv, w, r)
			} else if head == "auth" {
				auth(*sharedEnv, w, r)
			} else if head == "mockAuth" {
				mockAuth(*sharedEnv, w, r)
			} else {
				write.Error(errors.RouteNotFound, w, r)
			}
		default:
			write.Error(errors.BadRequestMethod, w, r)
		}
}



func getProfile(sharedEnv env.SharedEnv, u model.User, w http.ResponseWriter, r *http.Request) {
	walletAddress := r.URL.Query().Get("wallet_address")

	registeredApps, _ := sharedEnv.AppRepo().GetFromHost(walletAddress)
	playSessions, _ := sharedEnv.StreamSessionRepo().GetPlaySessions(walletAddress)
	hostSessions, _ := sharedEnv.StreamSessionRepo().GetHostSessions(walletAddress)

	isAllow := perm.RequireOwner(u.WalletAddress, walletAddress) &&
		perm.RequireAuthenticated(sharedEnv, w, r)
	if !isAllow {
		write.Error(errors.RouteUnauthorized, w, r)
		return
	}

	write.JSON(model.UserDetailProfileResponse{
		WalletAddress: walletAddress, CurUnreleasedBalance: 0, HourlyRate: 0,
		RegisteredApps: registeredApps, PlaySessions: playSessions, HostSessions: hostSessions,
	}, w, r)
}

func getUser(sharedEnv env.SharedEnv, w http.ResponseWriter, r *http.Request) {
	walletAddress := r.URL.Query().Get("wallet_address")

	dbUser, err := sharedEnv.UserRepo().GetUser(walletAddress)
	if err != nil {
		write.Error(err, w, r)
	}

	write.JSON(dbUser, w, r)
}

func signUp(sharedEnv env.SharedEnv, w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req model.SignUpReq
	err := decoder.Decode(&req)
	if err != nil || &req == nil {
		write.Error(errors.NoJSONBody, w, r)
	}

	dbUser, err := sharedEnv.UserRepo().SignUp(req.WalletAddress)
	if err != nil {
		write.Error(err, w, r)
	}

	write.JSON(dbUser, w, r)
}

func auth(sharedEnv env.SharedEnv, w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req model.AuthReq
	err := decoder.Decode(&req)
	if err != nil || &req == nil {
		write.Error(errors.NoJSONBody, w, r)
	}

	dbUser, err := sharedEnv.UserRepo().Auth(req.WalletAddress, req.Signature)
	if err != nil  {
		write.Error(err, w, r)
	}

	write.JSON(dbUser, w, r)
}

func mockAuth(sharedEnv env.SharedEnv, w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req model.MockAuthReq
	err := decoder.Decode(&req)
	if err != nil || &req == nil {
		write.Error(errors.NoJSONBody, w, r)
	}

	dbUser, _ := sharedEnv.UserRepo().GetUser(req.WalletAddress)
	jwt.WriteUserCookie(w, dbUser)
}