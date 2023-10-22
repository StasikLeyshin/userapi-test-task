package repository

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"refactoring/internal/models"
	"strconv"
	"sync"
	"time"
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

	userStore, err := u.getUserStore()
	if err != nil {
		return nil, err
	}

	return userStore.List, nil

}

func (u *UserRepository) CreateUser(user *models.CreateUserRequest) (*models.User, error) {
	u.mx.Lock()
	defer u.mx.Unlock()

	userStore, err := u.getUserStore()
	if err != nil {
		return nil, err
	}

	userStore.Increment++

	newUser := models.User{
		ID:          userStore.Increment,
		CreateDate:  time.Now(),
		DisplayName: user.DisplayName,
		Email:       user.Email,
	}

	id := strconv.Itoa(userStore.Increment)
	userStore.List[id] = newUser

	err = u.writeFile(*userStore)
	if err != nil {
		return nil, err
	}

	return &newUser, nil

}

func (u *UserRepository) GetUser(ID string) (*models.User, error) {

	userStore, err := u.getUserStore()
	if err != nil {
		return nil, err
	}

	user, ok := userStore.List[ID]
	if !ok {
		return nil, models.UserNotFound
	}
	return &user, nil
}

func (u *UserRepository) SearchUsers() (*models.UserStore, error) {
	var userStore models.UserStore

	file, err := os.ReadFile(u.filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening the file: %v", err)
	}

	err = json.Unmarshal(file, &userStore)
	if err != nil {
		return nil, fmt.Errorf("unmarshal failed: %v", err)
	}
	return &userStore, nil
}

func (u *UserRepository) UpdateUser(id string, newDisplayName string) (*models.User, error) {
	u.mx.Lock()
	defer u.mx.Unlock()

	userStore, err := u.getUserStore()
	if err != nil {
		return nil, err
	}

	user, ok := userStore.List[id]
	if !ok {
		return nil, models.UserNotFound
	}

	user.DisplayName = newDisplayName
	userStore.List[id] = user

	err = u.writeFile(*userStore)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepository) DeleteUser(id string) error {
	u.mx.Lock()
	defer u.mx.Unlock()

	userStore, err := u.getUserStore()
	if err != nil {
		return err
	}

	if _, ok := userStore.List[id]; !ok {
		return models.UserNotFound
	}

	delete(userStore.List, id)

	err = u.writeFile(*userStore)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserRepository) createFile() error {
	userStore := models.UserStore{
		Increment: 0,
		List:      map[string]models.User{},
	}

	err := u.writeFile(userStore)
	if err != nil {
		return err
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

func (u *UserRepository) getUserStore() (*models.UserStore, error) {
	var userStore models.UserStore
	file, err := os.ReadFile(u.filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening the file: %v", err)
	}

	err = json.Unmarshal(file, &userStore)
	if err != nil {
		return nil, fmt.Errorf("unmarshal failed:: %v", err)
	}
	return &userStore, nil
}

func (u *UserRepository) writeFile(userStore models.UserStore) error {

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
