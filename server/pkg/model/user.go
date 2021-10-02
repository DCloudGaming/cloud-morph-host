package model

import (
	"github.com/DCloudGaming/cloud-morph-host/pkg/utils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID string
	WalletAddress string `gorm:"primaryKey"`
	Nonce string `json:"nonce"`
	Machine string `json:"machine"`
	Location string `json:"location"`
	Name string `json:"name"`
	Status  Status
}

type SmartOtp struct {
	gorm.Model
	WalletAddress string `gorm:"primaryKey"`
	Otp string `json:"otp"`
}

type Status int

const (
	StatusDisabled Status = -1
	StatusUnverified = 0
	StatusActive = 1
	StatusAdmin = 10
)

type UserRepo interface {
	SignUp(walletAddress string) (*User, error)
	Auth(walletAddress string, signature string) (*User, error)
	GetUser(walletAddress string) (*User, error)
	UpdateUser(req UpdateUserReq) (*User, error)
	GenOTP(walletAddress string) (*SmartOtp, error)
	VerifyOTP(req VerifyOtpReq) (*SmartOtp, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepo {
	//db.AutoMigrate(&User{})
	return &userRepo{
		db,
	}
}

func (r *userRepo) SignUp(walletAddress string) (*User, error) {
	var nonce string
	nonce = utils.GenerateRandomString(10)
	user := User{WalletAddress: walletAddress, Nonce: nonce, Status: StatusUnverified}
	r.db.Create(&user)

	return &user, nil
}

func (r *userRepo) Auth(walletAddress string, signature string) (*User, error) {
	var user User
	r.db.First(&user, "wallet_address = ?", walletAddress)

	var msg string
	msg = "I am signing my one-time nonce: " + user.Nonce

	var verifyResult = utils.VerifySig(user.WalletAddress, signature, []byte(msg))

	if !verifyResult {
		return nil, errors.New("Wrong signature")
	}
	var newNonce string
	newNonce = utils.GenerateRandomString(10)
	r.db.Model(&user).Update("nonce", newNonce)
	return &user, nil
}

func (r *userRepo) GetUser(walletAddress string) (*User, error) {
	var user User
	r.db.First(&user, "wallet_address = ?", walletAddress)
	return &user, nil
}

func (r *userRepo) UpdateUser(req UpdateUserReq) (*User, error) {
	var user User
	r.db.First(&user, "wallet_address = ?", req.WalletAddress)
	user.Machine = req.Machine
	user.Location = req.Location
	user.Name = req.Name
	r.db.Save(&user)
	return &user, nil
}

func (r *userRepo) GenOTP(walletAddress string) (*SmartOtp, error) {
	otp := SmartOtp{WalletAddress: walletAddress, Otp: utils.GenerateRandomString(10)}
	r.db.Create(&otp)
	return &otp, nil
}

func (r *userRepo) VerifyOTP(req VerifyOtpReq) (*SmartOtp, error) {
	var smartOtp SmartOtp
	r.db.First(&smartOtp, "otp = ?", req.Otp)
	if &smartOtp != nil && smartOtp.Otp != req.Otp {
		return &smartOtp, nil
	}
	return &smartOtp, nil
}
