package storage

import "github.com/ignoxx/podara/poc3/types"

type MongoStorage struct{}

func (s *MongoStorage) Get(id string) *types.User {
	return &types.User{ID: id, Name: "Foo"}
}
