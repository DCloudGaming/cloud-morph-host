package env

import (
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
	"github.com/DCloudGaming/cloud-morph-host/pkg/model"
)

type sharedEnv struct {
	db *gorm.DB
	userRepo model.UserRepo
}

type SharedEnv interface {
	UserRepo() model.UserRepo
}

func New() (SharedEnv, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return &sharedEnv{
		db: db,
		userRepo: model.NewUserRepo(db),
	}, nil
}

func (e *sharedEnv) UserRepo() model.UserRepo {
	return e.userRepo
}
