package storage

import "github.com/ignoxx/podara/poc3/types"

type Storage interface {
	GetUserByEmail(string) (*types.User, error)
	CreateUser(string, string, string) (*types.User, error)
}
