package cache

import (
	"fmt"
	"refactoring/internal/models"
	"strconv"
	"sync"
)

type Cache struct {
	cache map[string]models.User
	mx    sync.Mutex
}

func NewCache() *Cache {
	return &Cache{cache: make(map[string]models.User)}
}

func (c *Cache) GetUser(id string) (*models.User, error) {
	user, ok := c.cache[id]
	fmt.Println(user, id)
	if !ok {
		return nil, models.UserNotFound
	}
	return &user, nil
}

func (c *Cache) AddUsers(userList models.UserList) error {
	c.mx.Lock()
	defer c.mx.Unlock()
	for key, value := range userList {
		c.cache[key] = value
	}
	return nil
}

func (c *Cache) AddUser(user models.User) error {
	c.mx.Lock()
	defer c.mx.Unlock()
	id := strconv.Itoa(user.ID)
	c.cache[id] = user
	return nil
}

func (c *Cache) UpdateUser(id string, newDisplayName string) error {
	c.mx.Lock()
	defer c.mx.Unlock()
	user, err := c.GetUser(id)
	if err != nil {
		return err
	}
	user.DisplayName = newDisplayName
	c.cache[id] = *user
	return nil
}
