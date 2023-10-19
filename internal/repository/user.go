package repository

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"refactoring/internal/models"
	"sync"
)

type UserRepository struct {
	logger   *log.Logger
	filePath string
	mx       sync.Mutex
}

func NewUserRepository(logger *log.Logger, filePath string) *UserRepository {
	return &UserRepository{
		logger:   logger,
		filePath: filePath,
	}
}

func (u *UserRepository) GetAllUsers() (models.UserList, error) {
	//var users []models.User
	var userStore models.UserStore

	file, err := os.ReadFile(u.filePath)

	if err != nil {
		return nil, fmt.Errorf("open file failed: %v", err)
	}

	err = json.Unmarshal(file, &userStore)
	if err != nil {
		return nil, fmt.Errorf("unmarshal failed: %v", err)
	}

	return userStore.List, nil

}

func (u *UserRepository) CreateUser() {
}

func (u *UserRepository) GetUser() {
}

func (u *UserRepository) SearchUsers() {
}

func (u *UserRepository) UpdateUser() {
}

func (u *UserRepository) DeleteUser() {
}

func (u *UserRepository) createFile() error {
	userStore := models.UserStore{
		Increment: 0,
		List:      map[string]models.User{},
	}
	data, err := json.Marshal(userStore)
	if err != nil {
		return fmt.Errorf("marshal failed: %v", err)
	}
	err = os.WriteFile(u.filePath, data, fs.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write data to file: %v", err)
	}
	return nil
}

func (u *UserRepository) CheckFileExist() error {
	_, err := os.Stat(u.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return u.createFile()
		}
		return fmt.Errorf("%v", err)
	}
	return nil
}
