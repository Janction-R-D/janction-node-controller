package repo

import (
	"node-controller/dao"
	model2 "node-controller/dao/model"
	"node-controller/model"
	"time"
)

type UserRepo struct {
}

func (u UserRepo) GetUserInfoByUserName(username string) (*model.User, error) {
	return dao.GetUserInfoByUserName(username)
}

func (u UserRepo) GetUserInfoByUserID(userID int) (*model.User, error) {
	return dao.GetUserInfoByUserID(userID)
}

func (u UserRepo) SetUserBlackToken(userID int, exp time.Time, token string) error {
	return dao.SetUserBlackToken(userID, exp, token)
}

func (u UserRepo) GetUserBlackToken(token string) (*model2.UserBlackToken, error) {
	return dao.GetUserBlackToken(token)
}

func (u UserRepo) CreateUser(params *model2.User) error {
	return dao.CreateUser(params)
}

func (u UserRepo) CreateUserInfo(params *model2.UserCreate) error {
	return dao.CreateUserInfo(params)
}

func (u UserRepo) GetAdminUserCount(search string) (int64, error) {
	return dao.GetAdminUserCount(search)
}

func (u UserRepo) GetAdminUserList(search string, limit, offset int) ([]*dao.UserInfo, error) {
	return dao.GetAdminUserList(search, limit, offset)
}

func (u UserRepo) UpdateUserPassed(id int, passwd string) error {
	return dao.UpdateUserPassed(id, passwd)
}

func (u UserRepo) DeleteUser(id int) error {
	return dao.DeleteUser(id)
}
