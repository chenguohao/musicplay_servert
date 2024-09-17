package dao

import "gorm.io/gorm"

type TUser struct {
	gorm.Model
	PostCount int
}

func (TUser) TableName() string {
	return "tuser"
}

type Post struct {
	gorm.Model
	UserId  uint
	Content string
}

func (Post) TableName() string {
	return "post"
}

func GetAll() ([]TUser, error) {
	var users []TUser
	//err := DBClient.Model(&TUser{}).Preload("Posts").Find(&users).Error
	err := DBClient.Select("tuser.*, (SELECT COUNT(*) FROM post WHERE post.user_id = tuser.id) as post_count").Find(&users).Error

	return users, err
}
