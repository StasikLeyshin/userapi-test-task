package cache

import (
	"errors"
	"refactoring/internal/models"
	"sync"
)

var (
	UserNotFound = errors.New("user_not_found")
)

type Cache struct {
	cache map[string]models.User
	mx    sync.Mutex
}

func NewCache() *Cache {
	return &Cache{cache: make(map[string]models.User)}
}

func (c *Cache) GetUser(ID string) (*models.User, error) {
	user, ok := c.cache[ID]
	if ok {
		return &user, nil
	}
	return nil, UserNotFound
}

func (c *Cache) AddUsers(userList models.UserList) error {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.cache = userList
	return nil
}

func (c *Cache) AddUser(ID string, user models.User) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.cache[ID] = user
}
