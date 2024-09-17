package dao

import "gorm.io/gorm"

type UserModel struct {
	gorm.Model
	Name        string `json:"name"`
	Uid         int    `json:"uid"`
	Avatar      string `json:"avatar"`
	Token       string `json:"token"`
	Email       string `json:"email"`
	PlatformUid string `json:"platform_uid"`
}

func (UserModel) TableName() string {
	return "user"
}

type SimpleUserModel struct {
	Name   string `json:"name"`
	Uid    int    `json:"uid"`
	Avatar string `json:"avatar"`
}

func (SimpleUserModel) TableName() string {
	return "user"
}

func AddUser(user UserModel) error {
	err := DBClient.Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

// 根据UID查询用户
func GetUserByUID(uid int) (*UserModel, error) {
	var user UserModel
	err := DBClient.Where("uid = ?", uid).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByPlatfromUid(platformID string) (*UserModel, error) {
	var user UserModel
	err := DBClient.Where("platform_uid = ?", platformID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// 计算用户数量
func GetUserCount() (int64, error) {
	var count int64
	err := DBClient.Model(&UserModel{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// 更新用户信息
func UpdateUserByUID(uid int, updatedUser UserModel) error {

	err := DBClient.Model(&UserModel{}).
		Where("uid = ?", uid).
		Select("name", "avatar"). // 只选择更新 name 和 avatar 字段
		Updates(map[string]interface{}{
			"name":   updatedUser.Name,
			"avatar": updatedUser.Avatar,
		}).Error
	if err != nil {
		return err
	}
	return nil

	//err := DBClient.Model(&UserModel{}).Where("uid = ?", uid).Updates(updatedUser).Error
	//if err != nil {
	//	return err
	//}
	//return nil
}
