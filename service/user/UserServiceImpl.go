package service

import (
	"errors"
	"fmt"

	"github.com/Tushar456/go-gin/entity"
)

type UserServiceImpl struct {
	users []*entity.User
}

func NewUserService() UserService {
	return &UserServiceImpl{
		users: []*entity.User{},
	}
}

func (u *UserServiceImpl) Save(user *entity.User) error {

	for _, us := range u.users {

		if us.Name == user.Name {
			return errors.New("user already exist")
		}

	}
	u.users = append(u.users, user)
	return nil
}

func (u *UserServiceImpl) Get(name *string) (*entity.User, error) {
	var findUser *entity.User

	for _, user := range u.users {

		if user.Name == *name {
			findUser = user
			return findUser, nil
		}

	}
	return nil, errors.New("user not found")

}

func (u *UserServiceImpl) GetAll() ([]*entity.User, error) {
	if len(u.users) == 0 {
		return nil, errors.New("no user found")
	}
	return u.users, nil

}

func (u *UserServiceImpl) Update(user *entity.User) error {
	for i, us := range u.users {

		if us.Name == user.Name {

			u.users[i] = user
			return nil
		}

	}
	return errors.New("user not found")
}

func (u *UserServiceImpl) Delete(name *string) error {
	for i, us := range u.users {

		if us.Name == *name {
			fmt.Println("inside")

			u.users = append(u.users[:i], u.users[i+1:]...)
			return nil
		}

	}
	return errors.New("user not found")

}
