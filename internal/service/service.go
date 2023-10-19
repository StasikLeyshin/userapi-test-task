package service

import (
	"log"
	"refactoring/internal/repository"
	"refactoring/internal/service/cache"
)

type Service struct {
	logger         *log.Logger
	userRepository *repository.UserRepository
	cache          *cache.Cache
}

func NewService(logger *log.Logger, filePath string) *Service {
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

func (s *Service) CreateUser() {
}

func (s *Service) GetUser(id int) {
}

func (s *Service) SearchUsers() {
}

func (s *Service) UpdateUser() {
}

func (s *Service) DeleteUser() {
}
