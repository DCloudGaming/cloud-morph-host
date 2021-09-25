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
}

type SharedEnv interface {
	UserRepo() model.UserRepo
	AppRepo() model.AppRepo
	HostConfigRepo() model.HostConfigRepo
	StreamSessionRepo() model.SessionRepo
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