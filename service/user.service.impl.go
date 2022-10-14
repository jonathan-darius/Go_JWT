package service

import (
	"Latihan_Mongo/model"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserServiceImpl struct {
	usercollection *mongo.Collection
	ctx            context.Context
}

func NewUserService(usercollection *mongo.Collection, ctx context.Context) UserService {
	return &UserServiceImpl{
		usercollection: usercollection,
		ctx:            ctx,
	}
}

func (u *UserServiceImpl) CreateUser(user *model.User_IN) error {
	_, err := u.usercollection.InsertOne(u.ctx, user)
	return err
}

func (u *UserServiceImpl) GetUser(name *string) (*model.User, error) {
	var user *model.User
	query := bson.D{
		bson.E{
			Key:   "user_name",
			Value: name,
		}}
	err := u.usercollection.FindOne(u.ctx, query).Decode(&user)
	return user, err
}

func (u *UserServiceImpl) GetAll() ([]*model.User, error) {
	var users []*model.User
	cursor, err := u.usercollection.Find(u.ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}
	for cursor.Next(u.ctx) {
		var user model.User
		err := cursor.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	cursor.Close(u.ctx)
	if len(users) == 0 {
		return nil, errors.New("No User Found")
	}
	return users, nil
}

func (u *UserServiceImpl) UpdateUser(user *model.User) error {
	filter := bson.D{bson.E{
		Key:   "user_name",
		Value: user.Name,
	}}
	update_data := bson.D{bson.E{Key: "$set", Value: bson.D{
		bson.E{Key: "user_name", Value: user.Name},
		bson.E{Key: "user_age", Value: user.Age},
		bson.E{Key: "user_address", Value: user.Address},
	}}}
	res, _ := u.usercollection.UpdateOne(u.ctx, filter, update_data)
	if res.MatchedCount != 1 {
		return errors.New("User Doesn't Match")
	}
	return nil
}

func (u *UserServiceImpl) DeleteUser(name *string) error {
	filter := bson.D{bson.E{
		Key:   "user_name",
		Value: name,
	}}
	res, _ := u.usercollection.DeleteOne(u.ctx, filter)
	if res.DeletedCount != 1 {
		return errors.New("User Doesn't Match")
	}
	return nil
}
