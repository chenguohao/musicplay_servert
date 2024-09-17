package services

import (
	"MusicPlayServer/dao"
	"fmt"
)

type LikeRequest struct {
	IsLike   bool  `json:"is_like"`
	TargetID int64 `json:"target_id"` // 假设这是一个数组，可以根据实际情况调整类型

}

func DoLike(req LikeRequest, curUserID int64) error {

	like := dao.LikeCountModel{}

	if req.IsLike {
		err := like.Add(curUserID, req.TargetID)
		if err != nil {
			fmt.Println("添加点赞失败:", err)
		} else {
			fmt.Println("点赞成功！")
		}
		return err
	} else {

		err := like.Cancel(curUserID, req.TargetID)
		return err
	}

	return nil

}
