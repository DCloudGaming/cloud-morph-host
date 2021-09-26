package model

import (
	"gorm.io/gorm"
)

type StreamSession struct {
	gorm.Model
	ID string `gorm:"autoIncrement"`
	StreamStatus StreamStatus `json:"stream_status"`
	MaxDuration int64 `json:"max_duration"`
	TotalDuration int64 `json:"total_duration"`
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
	GetSession(session_id int) (*StreamSession, error)
	StartSession(req StartSessionReq) (*StreamSession, error)
	UpdateSession(req UpdateSessionReq) (interface{}, error)
	GetPlaySessions(walletAddress string) ([]StreamSession, error)
	GetHostSessions(walletAddress string) ([]StreamSession, error)
}

type sessionRepo struct {
	db *gorm.DB
}

func NewSessionRepo(db *gorm.DB) SessionRepo {
	return &sessionRepo{
		db,
	}
}
func (r *sessionRepo) GetSession(session_id int) (*StreamSession, error) {
	var session StreamSession
	r.db.First(&session, "id = ?", session_id)
	return &session, nil
}

func (r *sessionRepo) GetPlaySessions(wallet_address string) ([]StreamSession, error) {
	var sessions []StreamSession
	dbRes := r.db.Find(&sessions, "client_wallet_address = ?", wallet_address)
	return sessions, dbRes.Error
}

func (r *sessionRepo) GetHostSessions(wallet_address string) ([]StreamSession, error) {
	var sessions []StreamSession
	dbRes := r.db.Find(&sessions, "host_wallet_address = ?", wallet_address)
	return sessions, dbRes.Error
}

func (r *sessionRepo) StartSession(req StartSessionReq) (*StreamSession, error) {
	session := StreamSession{
		StreamStatus: Streaming, MaxDuration: req.MaxDuration, AccumCharge: 0,
		ClientWalletAddress: req.ClientWalletAddress, HostWalletAddress: req.HostWalletAddress, AppName: req.AppName,
	}
	r.db.Create(&session)

	return &session, nil
}

func (r *sessionRepo) UpdateSession(req UpdateSessionReq) (interface{}, error) {
	var session StreamSession
	r.db.First(&session, "id = ?", req.SessionID)
	session.TotalDuration = req.TotalDuration
	session.AccumCharge = req.AccumCharge
	session.StreamStatus = StreamStatus(req.StreamStatus)
	r.db.Save(&session)

	return nil, nil
}
