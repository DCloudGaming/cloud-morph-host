package model

import (
	"gorm.io/gorm"
)

type HostConfig struct {
	gorm.Model
	ID string
	WalletAddress string `json:"wallet_address"`
	MaxConnections int `json:"max_connections"`
	CurUnreleasedBalance int64 `json:"cur_unreleased_balance"`
	HourlyRate int64 `json:"hourly_rate"`
}

type HostConfigRepo interface {

}

type hostConfigRepo struct {
	db *gorm.DB
}

func NewHostConfigRepo(db *gorm.DB) HostConfigRepo {
	return &hostConfigRepo{
		db,
	}
}
