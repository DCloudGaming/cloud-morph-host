package model

import (
	"github.com/DCloudGaming/cloud-morph-host/pkg/utils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type User struct {
	wallet_address string
	nonce string
}

type UserRepo interface {
	SignUp(walletAddress string) (*User, error)
	Auth(walletAddress string, signature string) (*User, error)
	GetUser(walletAddress string) (*User, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return &userRepo{
		db,
	}
}

func (r *userRepo) SignUp(walletAddress string) (*User, error) {
	var nonce string
	nonce = utils.GenerateRandomString(10)
	r.db.Create(&User{wallet_address: walletAddress, nonce: nonce})

	var user User
	r.db.First(&user, "walletAddress = ?", walletAddress)
	return &user, nil
}

func (r *userRepo) Auth(walletAddress string, signature string) (*User, error) {
	var user User
	r.db.First(&user, "walletAddress = ?", walletAddress)

	var msg string
	msg = "I am signing my one-time nonce: " + user.nonce

	var verifyResult = utils.VerifySig(user.wallet_address, signature, []byte(msg))

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
	r.db.First(&user, "walletAddress = ?", walletAddress)
	return &user, nil
}
