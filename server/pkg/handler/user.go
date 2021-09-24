package handler

import (
	"encoding/json"
	"fmt"
	"github.com/DCloudGaming/cloud-morph-host/pkg/env"
	"github.com/DCloudGaming/cloud-morph-host/pkg/errors"
	"github.com/DCloudGaming/cloud-morph-host/pkg/jwt"
	"github.com/DCloudGaming/cloud-morph-host/pkg/model"
	"github.com/DCloudGaming/cloud-morph-host/pkg/utils"
	"github.com/DCloudGaming/cloud-morph-host/pkg/write"
	"net/http"
)

func UserHandler(sharedEnv *env.SharedEnv) (f func(w http.ResponseWriter, r *http.Request)) {
	return func(w http.ResponseWriter, r *http.Request) {
		u, allowJwt := jwt.RequireAuth(model.StatusUnverified, *sharedEnv, w, r)
		if !allowJwt {
			return
		}
		var head string
		head, r.URL.Path = utils.ShiftPath(r.URL.Path)
		head, r.URL.Path = utils.ShiftPath(r.URL.Path)
		head, r.URL.Path = utils.ShiftPath(r.URL.Path)
		fmt.Println(u.WalletAddress)
		//var userDecode =
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
}



func getUser(sharedEnv env.SharedEnv, w http.ResponseWriter, r *http.Request) {
	walletAddress := r.URL.Query().Get("wallet_address")

	dbUser, err := sharedEnv.UserRepo().GetUser(walletAddress)
	if err != nil {
		write.Error(err, w, r)
	}

	write.JSON(dbUser, w, r)
}

type signUpReq struct {
	WalletAddress string `json:"wallet_address"`
}

func signUp(sharedEnv env.SharedEnv, w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req signUpReq
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

type authReq struct {
	WalletAddress string `json:"wallet_address"`
	Signature string `json:"signature"`
}

func auth(sharedEnv env.SharedEnv, w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req authReq
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

type mockAuthReq struct {
	WalletAddress string `json:"wallet_address"`
}

func mockAuth(sharedEnv env.SharedEnv, w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req mockAuthReq
	err := decoder.Decode(&req)
	if err != nil || &req == nil {
		write.Error(errors.NoJSONBody, w, r)
	}

	dbUser, _ := sharedEnv.UserRepo().GetUser(req.WalletAddress)
	jwt.WriteUserCookie(w, dbUser)
}