package storage

import (
	"errors"

	"github.com/google/uuid"
	"github.com/ignoxx/podara/poc3/types"
	"golang.org/x/crypto/bcrypt"
)

type MemoryStorage struct {
	users map[string]*types.User
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		users: make(map[string]*types.User),
	}
}

func (s *MemoryStorage) CreateUser(email, name, password string) (*types.User, error) {

	if s.users[email] != nil {
		return nil, errors.New("user already exists")
	}

	// hash password
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		return nil, err
	}

	user := &types.User{
		ID:       uuid.NewString(),
		Name:     name,
		Email:    email,
		Password: string(bytes),
	}

	s.users[email] = user

	return user, nil
}

func (s *MemoryStorage) GetUserByEmail(email string) (*types.User, error) {
	user := s.users[email]
	// check if user exists
	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}
