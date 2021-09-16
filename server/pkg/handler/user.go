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

func UserHandler(sharedEnv env.SharedEnv, w http.ResponseWriter, r *http.Request) http.HandlerFunc {
	var head string
	head, r.URL.Path = utils.ShiftPath(r.URL.Path)
	switch r.Method {
	case http.MethodGet:
		if head == "" {
			return getUser(sharedEnv, w, r)
		} else {
			return write.Error(errors.RouteNotFound)
		}
	case http.MethodPost:
		if head == "signup" {
			return signUp(sharedEnv, w, r)
		} else if head == "auth" {
			return auth(sharedEnv, w, r)
		} else {
			return write.Error(errors.RouteNotFound)
		}
	default:
		return write.Error(errors.BadRequestMethod)
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
		return write.Error(errors.NoJSONBody)
	}

	dbUser, err := sharedEnv.UserRepo().GetUser(req.walletAddress)
	if err != nil {
		return write.Error(err)
	}

	return write.JSON(dbUser)
}

type signUpReq struct {
	walletAddress string
}

func signUp(sharedEnv env.SharedEnv, w http.ResponseWriter, r *http.Request) http.HandlerFunc {
	decoder := json.NewDecoder(r.Body)
	var req signUpReq
	err := decoder.Decode(&req)
	if err != nil || &req == nil {
		return write.Error(errors.NoJSONBody)
	}

	dbUser, err := sharedEnv.UserRepo().SignUp(req.walletAddress)
	if err != nil {
		return write.Error(err)
	}

	return write.JSON(dbUser)
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
		return write.Error(errors.NoJSONBody)
	}

	dbUser, err := sharedEnv.UserRepo().Auth(req.walletAddress, req.signature)
	if err != nil  {
		return write.Error(err)
	}

	return write.JSON(dbUser)
}

