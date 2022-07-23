package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"

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

func (u *UserServiceImpl) Save(ctx *gin.Context, user *entity.User) error {

	for _, us := range u.users {

		if us.UserName == user.UserName {
			return errors.New("user already exist")
		}

	}
	u.users = append(u.users, user)
	return nil
}

func (u *UserServiceImpl) Get(ctx *gin.Context, username *string) (*entity.User, error) {
	var findUser *entity.User

	for _, user := range u.users {

		if user.UserName == *username {
			findUser = user
			return findUser, nil
		}

	}
	return nil, errors.New("user not found")

}

func (u *UserServiceImpl) GetAll(ctx *gin.Context) ([]*entity.User, error) {
	if len(u.users) == 0 {
		return nil, errors.New("no user found")
	}
	return u.users, nil

}

func (u *UserServiceImpl) Update(ctx *gin.Context, user *entity.User) error {
	for i, us := range u.users {

		if us.UserName == user.UserName {

			u.users[i] = user
			return nil
		}

	}
	return errors.New("user not found")
}

func (u *UserServiceImpl) Delete(ctx *gin.Context, username *string) error {
	for i, us := range u.users {

		if us.UserName == *username {
			fmt.Println("inside")

			u.users = append(u.users[:i], u.users[i+1:]...)
			return nil
		}

	}
	return errors.New("user not found")

}
