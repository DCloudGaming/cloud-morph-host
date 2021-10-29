package env

import (
	"github.com/DCloudGaming/cloud-morph-host/pkg/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type sharedEnv struct {
	db *gorm.DB
	userRepo model.UserRepo
	appRepo model.AppRepo
	hostConfigRepo model.HostConfigRepo
	streamSessionRepo model.SessionRepo
	mode string // DEBUG or PROD
	defaultAppPath string
}

type SharedEnv interface {
	UserRepo() model.UserRepo
	AppRepo() model.AppRepo
	HostConfigRepo() model.HostConfigRepo
	StreamSessionRepo() model.SessionRepo
	Mode() string
	DefaultAppPath() string
}

func New() (SharedEnv, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return &sharedEnv{
		db: db,
		userRepo: model.NewUserRepo(db),
		appRepo: model.NewAppRepo(db),
		hostConfigRepo: model.NewHostConfigRepo(db),
		streamSessionRepo: model.NewSessionRepo(db),
		//mode: "DEBUG",
		mode: "PROD",
		defaultAppPath: "/Users/hieuletrung/Documents/repos/side_projects/cloud-morph-host/streamer/apps/Minesweeper.exe",
	}, nil
}

func (e *sharedEnv) UserRepo() model.UserRepo {
	return e.userRepo
}

func (e *sharedEnv) AppRepo() model.AppRepo {
	return e.appRepo
}

func (e *sharedEnv) HostConfigRepo() model.HostConfigRepo {
	return e.hostConfigRepo
}

func (e *sharedEnv) StreamSessionRepo() model.SessionRepo {
	return e.streamSessionRepo
}

func (e *sharedEnv) Mode() string {
	return e.mode
}

func (e *sharedEnv) DefaultAppPath() string {
	return e.defaultAppPath
}
