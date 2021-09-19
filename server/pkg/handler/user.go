package handler

import (
	"encoding/json"
	"net/http"
	"github.com/miguelmota/go-ethutil"

	"github.com/DCloudGaming/cloud-morph-host/pkg/env"
	"github.com/DCloudGaming/cloud-morph-host/pkg/errors"
	"github.com/DCloudGaming/cloud-morph-host/pkg/model"
	"github.com/DCloudGaming/cloud-morph-host/pkg/utils"
	"github.com/DCloudGaming/cloud-morph-host/pkg/write"
)

func UserHandler(sharedEnv env.SharedEnv, w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = utils.ShiftPath(r.URL.Path)
	switch r.Method {
	case http.MethodGet:
		if head == "" {
			getUser(sharedEnv, w, r)
		} else {
			write.Error(errors.RouteNotFound, w, r)
		}
	case http.MethodPost:
		if head == "signup" {
			signUp(sharedEnv, w, r)
		} else if head == "auth" {
			auth(sharedEnv, w, r)
		} else {
			write.Error(errors.RouteNotFound, w, r)
		}
	default:
		write.Error(errors.BadRequestMethod, w, r)
	}
}

type getUserReq struct {
	walletAddress string
}

func getUser(sharedEnv env.SharedEnv, w http.ResponseWriter, r *http.Request) http.HandlerFunc {
	decoder := json.NewDecoder(r.Body)
	var req getUserReq
	err := decoder.Decode(&req)
	if err != nil || &req == nil {
		write.Error(errors.NoJSONBody, w, r)
	}

	dbUser, err := sharedEnv.UserRepo().GetUser(req.walletAddress)
	if err != nil {
		write.Error(err, w, r)
	}

	write.JSON(dbUser, w, r)
}

type signUpReq struct {
	walletAddress string
}

func signUp(sharedEnv env.SharedEnv, w http.ResponseWriter, r *http.Request) http.HandlerFunc {
	decoder := json.NewDecoder(r.Body)
	var req signUpReq
	err := decoder.Decode(&req)
	if err != nil || &req == nil {
		write.Error(errors.NoJSONBody, w, r)
	}

	dbUser, err := sharedEnv.UserRepo().SignUp(req.walletAddress)
	if err != nil {
		write.Error(err, w, r)
	}

	write.JSON(dbUser, w, r)
}

type authReq struct {
	walletAddress string
	signature string
}

func auth(sharedEnv env.SharedEnv, w http.ResponseWriter, r *http.Request) http.HandlerFunc {
	decoder := json.NewDecoder(r.Body)
	var req authReq
	err := decoder.Decode(&req)
	if err != nil || &req == nil {
		write.Error(errors.NoJSONBody, w, r)
	}

	dbUser, err := sharedEnv.UserRepo().Auth(req.walletAddress, req.signature)
	if err != nil  {
		write.Error(err, w, r)
	}

	write.JSON(dbUser, w, r)
}

