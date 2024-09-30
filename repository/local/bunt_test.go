package local

import (
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/model"
	"go.uber.org/zap"
	"os"
	"testing"
)

func TestBuntRepositoryCreate(t *testing.T) {
	log, _ := zap.NewDevelopment()

	repo, err := NewBuntRepository(os.Getenv("LOCAL_DB_PATH"))
	if err != nil {
		log.Error("Error creating brand-new bunt repository", zap.Error(err))
		t.FailNow()
	}

	user := model.User{
		ID:       "1",
		Username: "",
		Email:    "",
		Password: "",
		Name:     "",
		Lastname: "",
		Age:      0,
	}

	if err = repo.Create(&user); err != nil {
		log.Error("Error creating new user", zap.Any("user", user), zap.Error(err))
		t.FailNow()
	}

	log.Info("User created successfully")

}
