package service

import (
	"Latihan_Mongo/model"
)

type UserService interface {
	CreateUser(*model.User_IN) error
	GetUser(*string) (*model.User, error)
	GetAll() ([]*model.User, error)
	UpdateUser(*model.User) error
	DeleteUser(*string) error
}
