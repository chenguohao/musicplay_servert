package services

import (
	"MusicPlayServer/dao"
	"encoding/json"
	"errors"
	"fmt"
)

type CreatePlaylistRequest struct {
	Title    string        `json:"title" binding:"required"`
	CoverURL string        `json:"cover" binding:"required"`
	ListItem []interface{} `json:"playlist" binding:"required"` // 假设这是一个数组，可以根据实际情况调整类型
	OwnerID  int           `json:"ownerID" binding:"required"`
}

func CreatePlaylist(req CreatePlaylistRequest) error {

	playlist := dao.PlaylistModel{}

	playlist.Title = req.Title

	playlist.CoverURL = req.CoverURL
	listItemJSON, err := json.Marshal(req.ListItem)
	if err != nil {
		fmt.Println("Failed to marshal ListItem:", err)
		return fmt.Errorf("error code 10011: %w", errors.New("param error"))
	}
	playlist.ListItem = listItemJSON
	//playlist.ListItem = req.ListItem
	playlist.PlaylistID = generateUUID() // 你可以使用任何合适的 UUID 生成方法

	// 调用 AddPlaylist 将数据存入数据库
	err = dao.AddPlaylist(playlist)
	return err
}

func generateUUID() int64 {
	curCount, _ := dao.GetPlaylistCount()
	return 21000 + curCount
}

type GetPlaylistRequest struct {
	Page int `json:"page" binding:"required"`
	Size int `json:"size" binding:"required"`
}

func GetPlayList(page int, size int) ([]dao.PlaylistModel, error) {
	return dao.GetPlaylistsByPage(page, size)
}
