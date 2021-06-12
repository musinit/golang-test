package repository_test

import (
	"golang-test/db"
	"golang-test/repository"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/google/uuid"
)

func Test_GetUserByID_Success(t *testing.T) {
	db, _ := db.SetDB()
	userRepository := repository.UserRepository{db}

	_, err := userRepository.GetAllUsers()

	assert.Nil(t, err)
}

func Test_Create_Success(t *testing.T) {
	db, _ := db.SetDB()
	userRepository := repository.UserRepository{db}

	testID, _ := uuid.Parse("daa36a15-5810-477e-b97a-309c9b78275f")
	err := userRepository.Create(&repository.User{
		ID:         testID,
		Name:       "Александр",
		Surname:    "Бугай",
		Patronymic: "Не помню, прости",
	})

	assert.Nil(t, err)
}
