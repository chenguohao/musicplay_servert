package dao

import "gorm.io/gorm"

type PlayCountModel struct {
	gorm.Model
	Sender int64 `json:"sender"`
	Target int64 `json:"target"`
}

func (PlayCountModel) TableName() string {
	return "playcount"
}

func (model *PlayCountModel) Add(sender int64, target int64) (bool, error) {
	var play PlayCountModel

	// 检查是否已经存在 sender-target 组合
	err := DBClient.Where("sender = ? AND target = ?", sender, target).
		First(&play).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 如果不存在，则创建新的记录
			newPlay := PlayCountModel{
				Sender: sender,
				Target: target,
			}
			if err := DBClient.Create(&newPlay).Error; err != nil {
				return false, err
			} else {
				return true, nil
			}
		}
	} else {
		// 查询错误
		return false, err
	}
	return false, nil
}
