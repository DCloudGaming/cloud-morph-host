package handler

import (
	"encoding/json"
	"github.com/DCloudGaming/cloud-morph-host/pkg/env"
	"github.com/DCloudGaming/cloud-morph-host/pkg/errors"
	"github.com/DCloudGaming/cloud-morph-host/pkg/jwt"
	"github.com/DCloudGaming/cloud-morph-host/pkg/model"
	"github.com/DCloudGaming/cloud-morph-host/pkg/perm"
	"github.com/DCloudGaming/cloud-morph-host/pkg/write"
	"gorm.io/gorm"
	"net/http"
)

func UserHandler(
	sharedEnv *env.SharedEnv, w http.ResponseWriter, r *http.Request,
	u *model.User, hostU *model.User, head string) {
	switch r.Method {
		case http.MethodGet:
			if head == "" {
				getUser(*sharedEnv, w, r)
			} else if head == "getFromToken" {
				getUserFromToken(*sharedEnv, *u, w, r)
			} else if head == "genOTP" {
				genOTP(*sharedEnv, *u, w, r)
			} else if head == "profile" {
				getProfile(*sharedEnv, *u, w, r)
			} else if head == "getAdminSettings" {
				getAdminSettings(*sharedEnv, *u, w, r)
			} else {
				write.Error(errors.RouteNotFound, w, r)
			}
		case http.MethodPost:
			if head == "signup" {
				signUp(*sharedEnv, w, r)
			} else if head == "getOrCreate" {
				getOrCreate(*sharedEnv, *u, w, r)
			} else if head == "auth" {
				auth(*sharedEnv, w, r)
			} else if head == "mockAuth" {
				mockAuth(*sharedEnv, w, r)
			} else if head == "update" {
				updateUser(*sharedEnv, *u, w, r)
			} else if head == "updateAdminSettings" {
				updateAdminSettings(*sharedEnv, *u, w, r)
			} else if head == "verifyOTP" {
				verifyOTP(*sharedEnv, w, r)
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
	dbUser, _ := sharedEnv.UserRepo().GetUser(walletAddress)


	isAllow := perm.RequireOwner(sharedEnv, u.WalletAddress, walletAddress) &&
		perm.RequireAuthenticated(sharedEnv, w, r)
	if !isAllow {
		write.Error(errors.RouteUnauthorized, w, r)
		return
	}

	write.JSON(model.UserDetailProfileResponse{
		WalletAddress: walletAddress, CurUnreleasedBalance: 0, HourlyRate: 0, Location: dbUser.Location, Machine: dbUser.Machine,
		RegisteredApps: registeredApps, PlaySessions: playSessions, HostSessions: hostSessions,
	}, w, r)
}

func getUserFromToken(sharedEnv env.SharedEnv, u model.User, w http.ResponseWriter, r *http.Request) {

	dbUser, err := sharedEnv.UserRepo().GetUser(u.WalletAddress)
	if err != nil {
		write.Error(err, w, r)
	}

	write.JSON(dbUser, w, r)
}

func getAdminSettings(sharedEnv env.SharedEnv, u model.User, w http.ResponseWriter, r *http.Request) {
	dbAdminConfigs, _ := sharedEnv.UserRepo().GetAdminSettings()
	allowApps, _ := sharedEnv.AppRepo().GetAllowedApps()

	var resp model.GetAdminConfigsResponse
	resp.HourlyRate = dbAdminConfigs.HourlyRate
	resp.AllowedApps = []model.AllowAppSchema{}
	for _, allowApp := range allowApps {
		resp.AllowedApps = append(resp.AllowedApps, model.AllowAppSchema{
			AppName: allowApp.AppName, ImageUrl: allowApp.ImageUrl, Publisher: allowApp.Publisher,
		})
	}

	write.JSON(resp, w, r)
}

func updateAdminSettings(sharedEnv env.SharedEnv, u model.User, w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req model.UpdateAdminReq
	err := decoder.Decode(&req)
	if err != nil || &req == nil {
		write.Error(errors.NoJSONBody, w, r)
		return
	}

	isAllow := perm.RequireAdmin(sharedEnv, u.WalletAddress)
	if !isAllow {
		write.Error(errors.RouteUnauthorized, w, r)
		return
	}

	sharedEnv.UserRepo().UpdateAdminSettings(req)
	sharedEnv.AppRepo().RemoveUnallowedAppsFromRegister(req.AllowedApps)
	sharedEnv.AppRepo().AllowNewApps(req.AllowedApps)
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

func updateUser(sharedEnv env.SharedEnv, u model.User, w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req model.UpdateUserReq
	err := decoder.Decode(&req)
	if err != nil || &req == nil {
		write.Error(errors.NoJSONBody, w, r)
		return
	}

	isAllow := perm.RequireOwner(sharedEnv, u.WalletAddress, req.WalletAddress) &&
		perm.RequireAuthenticated(sharedEnv, w, r)
	if !isAllow {
		write.Error(errors.RouteUnauthorized, w, r)
		return
	}

	dbUser, _ := sharedEnv.UserRepo().UpdateUser(req)
	write.JSON(dbUser, w, r)
}

// Also validate
func getOrCreate(sharedEnv env.SharedEnv, u model.User, w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req model.GetOrCreateUserReq
	err := decoder.Decode(&req)
	if err != nil || &req == nil {
		write.Error(errors.NoJSONBody, w, r)
		return
	}

	if (u.WalletAddress == req.WalletAddress) {
		return
	}

	dbUser, err := sharedEnv.UserRepo().GetUser(req.WalletAddress)
	if u.WalletAddress != dbUser.WalletAddress && err == nil {
		write.JSON(dbUser, w, r)
		return
	} else if err != nil && err != gorm.ErrRecordNotFound  {
		write.Error(err, w, r)
		return
	} else if err == gorm.ErrRecordNotFound {
		dbUser2, err2 := sharedEnv.UserRepo().SignUp(req.WalletAddress)
		if err2 != nil {
			write.Error(err2, w, r)
		}
		write.JSON(dbUser2, w, r)
		return
	}

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
		return
	}
	jwt.WriteUserCookie(w, dbUser)
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

func genOTP(sharedEnv env.SharedEnv, u model.User, w http.ResponseWriter, r *http.Request) {
	otp, _ := sharedEnv.UserRepo().GenOTP(u.WalletAddress)
	write.JSON(otp, w, r)
}

func verifyOTP(sharedEnv env.SharedEnv, w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req model.VerifyOtpReq
	err := decoder.Decode(&req)
	if err != nil || &req == nil {
		write.Error(errors.NoJSONBody, w, r)
	}
	smartOtp, _ := sharedEnv.UserRepo().VerifyOTP(req)
	if smartOtp != nil {
		dbUser, _ := sharedEnv.UserRepo().GetUser(smartOtp.WalletAddress)
		var resp model.VerifyOTPResponse
		resp.WalletAddress = dbUser.WalletAddress
		resp.Token = jwt.EncodeUser(dbUser)
		jwt.WriteUserCookie(w, dbUser)
		write.JSON(resp, w, r)
	}
}
