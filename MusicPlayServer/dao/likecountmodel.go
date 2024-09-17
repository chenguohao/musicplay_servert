package dao

import "gorm.io/gorm"

type LikeCountModel struct {
	gorm.Model
	Sender   int64 `json:"sender"`
	Target   int64 `json:"target"`
	IsEnable bool  `json:"is_enable"`
}

func (LikeCountModel) TableName() string {
	return "likecount"
}

func (model *LikeCountModel) Count(target int64) (int64, error) {
	var count int64
	err := DBClient.Model(&LikeCountModel{}).
		Where("target = ? AND is_enable = ?", target, true).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (model *LikeCountModel) Add(sender int64, target int64) error {
	var like LikeCountModel

	// 检查是否已经存在 sender-target 组合
	err := DBClient.Where("sender = ? AND target = ?", sender, target).
		First(&like).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 如果不存在，则创建新的记录
			newLike := LikeCountModel{
				Sender:   sender,
				Target:   target,
				IsEnable: true,
			}
			if err := DBClient.Create(&newLike).Error; err != nil {
				return err
			}
		} else {
			// 查询错误
			return err
		}
	} else {
		// 如果已存在记录
		if like.IsEnable == false {
			// 如果 is_enable 为 false，则更新为 true
			like.IsEnable = true
			if err := DBClient.Save(&like).Error; err != nil {
				return err
			}
		}
		// 如果 is_enable 为 true，则什么都不做
	}
	return nil
}

func (model *LikeCountModel) Cancel(sender int64, target int64) error {
	var like LikeCountModel

	// 查找对应的记录
	err := DBClient.Where("sender = ? AND target = ?", sender, target).
		First(&like).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 如果没有找到记录，直接返回，不执行任何操作
			return nil
		}
		// 查询错误
		return err
	}

	// 如果找到记录且 is_enable 为 true，则将其更新为 false
	if like.IsEnable == true {
		like.IsEnable = false
		if err := DBClient.Save(&like).Error; err != nil {
			return err
		}
	} else {
		like.IsEnable = true
		if err := DBClient.Save(&like).Error; err != nil {
			return err
		}
	}
	// 如果 is_enable 已经是 false，不做任何操作
	return nil
}
