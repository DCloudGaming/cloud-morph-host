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

type AllowedApp struct {
	gorm.Model
	ID string
	AppName string `json:"app_name"`
	Publisher string `json:"publisher"`
	ImageUrl string `json:"image_url"`
}

type AllowAppSchema struct {
	AppName string `json:"app_name"`
	ImageUrl string `json:"image_url"`
	Publisher string `json:"publisher"`
}

type AppVote struct {
	gorm.Model
	ID string
	AppName string `json:"app_name"`
	WalletAddress string `json:"wallet_address"`
}

type AppRepo interface {
	RegisterBatch(req RegisterAppReq) (int64, error)
	GetFromHost(walletAddress string) ([]RegisteredApp, error)
	GetAppByName(appName string, walletAddress string) (RegisteredApp, error)
	GetAllRegisteredApps() ([]RegisteredApp, error)
	RemoveUnallowedAppsFromRegister(apps []AllowAppSchema) ()
	AllowNewApps(apps []AllowAppSchema) ()
	DisallowApps(appNames []string) ()
	GetAllowedApps() ([]AllowedApp, error)
	UpdateVote(appName string, walletAddress string) ()
	GetVote(appName string) (int)
	IsVoted(appName string, walletAddress string) (bool)
}

type appRepo struct {
	db *gorm.DB
}

func NewAppRepo(db *gorm.DB) AppRepo {
	return &appRepo{
		db,
	}
}

func (r *appRepo) RegisterBatch(req RegisterAppReq) (int64, error) {
	r.db.Where("wallet_address = ?", req.WalletAddress).Unscoped().Delete(&RegisteredApp{})

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

func (r *appRepo) GetAppByName(appName string, walletAddress string) (RegisteredApp, error) {
	var registeredApp RegisteredApp
	dbRes := r.db.First(&registeredApp, "wallet_address = ? and app_name = ?", walletAddress, appName)
	return registeredApp, dbRes.Error
}

func (r *appRepo) GetAllRegisteredApps() ([]RegisteredApp, error) {
	var apps []RegisteredApp
	dbRes := r.db.Find(&apps)
	return apps, dbRes.Error
}

func (r *appRepo) RemoveUnallowedAppsFromRegister(apps []AllowAppSchema) () {
	var appNames []string
	for i := 0; i < len(apps); i ++ {
		appNames = append(appNames, apps[i].AppName)
	}
	r.db.Where("app_name NOT IN ?", appNames).Unscoped().Delete(&RegisteredApp{})
}

func (r *appRepo) AllowNewApps(apps []AllowAppSchema) () {
	r.db.Where("1=1").Unscoped().Delete(&AllowedApp{})
	var allowApps = []AllowedApp{}
	for i := 0; i < len(apps); i ++ {
		allowApps = append(allowApps, AllowedApp{
			AppName: apps[i].AppName, Publisher: apps[i].Publisher, ImageUrl: apps[i].ImageUrl,
		})
	}
	r.db.Create(&allowApps)
}

func (r *appRepo) DisallowApps(appNames []string) () {
	r.db.Where("app_name IN ?" , appNames).Unscoped().Delete(&AllowedApp{})
}

func (r *appRepo) GetAllowedApps() ([]AllowedApp, error) {
	var allowApps []AllowedApp
	dbRes := r.db.Find(&allowApps)
	return allowApps, dbRes.Error
}

func (r *appRepo) UpdateVote(appName string, walletAddress string) {
	var vote AppVote
	err := r.db.First(&vote, "wallet_address = ? and app_name = ?", walletAddress, appName).Error
	if (err != gorm.ErrRecordNotFound) {
		r.db.Where("wallet_address = ? and app_name = ?", walletAddress, appName).Unscoped().Delete(&vote)
	} else {
		r.db.Create(&AppVote{WalletAddress: walletAddress, AppName: appName})
	}
}

func (r *appRepo) IsVoted(appName string, walletAddress string) (bool) {
	var vote AppVote
	err := r.db.First(&vote, "wallet_address = ? and app_name = ?", walletAddress, appName).Error
	return err == nil
}

func (r *appRepo) GetVote(appName string) (int) {
	type result struct {
		Num int
	}
	var n result
	r.db.Model(&AppVote{}).Select("count(*) as num").Where("app_name = ?", appName).Scan(&n)
	if (&n != nil) {
		return n.Num
	} else {
		return 0
	}
}
