package model

import (
	"gorm.io/gorm"
)

type StreamSession struct {
	gorm.Model
	ID string
	StreamStatus StreamStatus `json:"stream_status"`
	MaxDuration int64 `json:"max_duration"`
	AccumCharge int64 `json:"accum_charge"`
	ClientWalletAddress string `json:"client_wallet_address"`
	HostWalletAddress string `json:"host_wallet_address"`
	AppName string `json:"app_name"`
}

type StreamStatus int
const (
	Streaming StreamStatus = 0
	Paused = 1
	Finished = 2
)

type SessionRepo interface {

}

type sessionRepo struct {
	db *gorm.DB
}

func NewSessionRepo(db *gorm.DB) SessionRepo {
	return &sessionRepo{
		db,
	}
}
