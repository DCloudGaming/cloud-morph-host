package model

import (
	"gorm.io/gorm"
)

type RegisteredApp struct {
	gorm.Model
	ID string
	WalletAddress string `json:"wallet_address"`
	AppPath string `json:"app_path"`
	AppName string `json:"app_name"`
}

type AppRepo interface {
	RegisterBatch(req RegisterAppReq) (int64, error)
	GetFromHost(walletAddress string) ([]RegisteredApp, error)
	GetAllRegisteredApps() ([]RegisteredApp, error)
}

type appRepo struct {
	db *gorm.DB
}

func NewAppRepo(db *gorm.DB) AppRepo {
	return &appRepo{
		db,
	}
}

// TODO: Only register the unregistered apps, and ignore the rest. For now it might fail
func (r *appRepo) RegisterBatch(req RegisterAppReq) (int64, error) {
	var apps = []RegisteredApp{}
	// TODO: Refactor
	for i := 0; i < len(req.AppPaths); i++ {
		apps = append(apps, RegisteredApp{WalletAddress: req.WalletAddress, AppPath: req.AppPaths[i], AppName: req.AppNames[i]})
	}
	dbRes := r.db.Create(&apps)

	return dbRes.RowsAffected, dbRes.Error
}

func (r *appRepo) GetFromHost(walletAddress string) ([]RegisteredApp, error) {
	var apps []RegisteredApp
	dbRes := r.db.Find(&apps, "wallet_address = ?", walletAddress)
	return apps, dbRes.Error
}

func (r *appRepo) GetAllRegisteredApps() ([]RegisteredApp, error) {
	var apps []RegisteredApp
	dbRes := r.db.Find(&apps)
	return apps, dbRes.Error
}
