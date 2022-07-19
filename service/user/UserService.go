package service

import "github.com/Tushar456/go-gin/entity"

type UserService interface {
	Save(*entity.User) error
	Get(*string) (*entity.User, error)
	GetAll() ([]*entity.User, error)
	Update(*entity.User) error
	Delete(*string) error
}
