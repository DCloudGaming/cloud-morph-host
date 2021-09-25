package handler

import (
	"encoding/json"
	"github.com/DCloudGaming/cloud-morph-host/pkg/env"
	"github.com/DCloudGaming/cloud-morph-host/pkg/errors"
	"github.com/DCloudGaming/cloud-morph-host/pkg/jwt"
	"github.com/DCloudGaming/cloud-morph-host/pkg/model"
	"github.com/DCloudGaming/cloud-morph-host/pkg/schema"
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
	var req schema.SignUpReq
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
	var req schema.AuthReq
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
	var req schema.MockAuthReq
	err := decoder.Decode(&req)
	if err != nil || &req == nil {
		write.Error(errors.NoJSONBody, w, r)
	}

	dbUser, _ := sharedEnv.UserRepo().GetUser(req.WalletAddress)
	jwt.WriteUserCookie(w, dbUser)
}