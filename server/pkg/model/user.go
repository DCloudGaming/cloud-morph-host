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
	Status  Status
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
