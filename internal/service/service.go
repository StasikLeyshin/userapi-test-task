package service

import (
	"github.com/sirupsen/logrus"
	"refactoring/internal/models"
	"refactoring/internal/repository"
	"refactoring/internal/service/cache"
)

type Service struct {
	logger         *logrus.Logger
	userRepository *repository.UserRepository
	cache          *cache.Cache
}

func NewService(logger *logrus.Logger, filePath string) *Service {
	return &Service{
		logger:         logger,
		userRepository: repository.NewUserRepository(logger, filePath),
		cache:          cache.NewCache(),
	}

}

func (s *Service) Start() error {
	err := s.userRepository.CheckFileExist()
	if err != nil {
		return err
	}

	users, err := s.userRepository.GetAllUsers()
	if err != nil {
		return err
	}

	err = s.cache.AddUsers(users)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) CreateUser(user *models.CreateUserRequest) (*models.User, error) {
	createUser, err := s.userRepository.CreateUser(user)
	if err != nil {
		return nil, err
	}

	_ = s.cache.AddUser(*createUser)

	return createUser, nil
}

func (s *Service) GetUser(id string) (*models.User, error) {
	user, err := s.cache.GetUser(id)
	if err == nil {
		return user, nil
	}

	user, err = s.userRepository.GetUser(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) SearchUsers() (*models.UserStore, error) {
	users, err := s.userRepository.SearchUsers()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *Service) UpdateUser(id string, newDisplayName string) (*models.User, error) {
	user, err := s.userRepository.UpdateUser(id, newDisplayName)
	if err != nil {
		return nil, err
	}
	err = s.cache.UpdateUser(id, newDisplayName)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) DeleteUser(id string) error {
	err := s.userRepository.DeleteUser(id)
	if err != nil {
		return err
	}

	err = s.cache.DeleteUser(id)
	if err != nil {
		return err
	}

	return nil
}
