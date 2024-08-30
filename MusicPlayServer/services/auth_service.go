package services

import (
	"MusicPlayServer/dao"
)

type AppleLoginRequest struct {
	AuthorizationCode string `json:"authorizationCode" binding:"required"`
	UserID            string `json:"userID" binding:"required"`
	Email             string `json:"email"`
	FullName          string `json:"fullName"`
	IDToken           string `json:"idToken" binding:"required"`
}

func RegisterOrLogin(req AppleLoginRequest) dao.UserModel {

	user, _ := dao.GetUserByPlatfromUid(req.UserID)
	if user != nil {
		return *user
	}

	newUser := dao.UserModel{}
	newUser.PlatformUid = req.UserID
	newUser.Email = req.Email
	newUser.Name = req.FullName
	newUser.Token = req.IDToken

	count, _ := dao.GetUserCount()
	newUser.Uid = int(10000 + count)

	dao.AddUser(newUser)

	return newUser
}
