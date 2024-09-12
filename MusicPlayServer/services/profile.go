package services

import "MusicPlayServer/dao"

func UpdateProfile(req ProfileRequest) error {

	user := dao.UserModel{}
	user.Name = req.Name
	user.Avatar = req.Avatar
	return dao.UpdateUserByUID(req.UserID, user)
}
