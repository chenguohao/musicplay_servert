package services

import (
	"MusicPlayServer/dao"
	"encoding/json"
	"errors"
	"fmt"
)

type CreatePlaylistRequest struct {
	Title    string        `json:"title" binding:"required"`
	CoverURL string        `json:"cover_url"`
	ListItem []interface{} `json:"list_item" binding:"required"` // 假设这是一个数组，可以根据实际情况调整类型
	OwnerID  int64         `json:"owner_id" binding:"required"`
}

type DeletePlaylistRequest struct {
	PlaylistID int64 `json:"playlist_id" binding:"required"`
}

func CreatePlaylist(req CreatePlaylistRequest) error {

	playlist := dao.PlaylistModel{}

	playlist.Title = req.Title
	playlist.OwnerID = req.OwnerID
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
	err = dao.AddNewPlaylist(playlist)
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

func GetPlayList(page int, size int, curUid int64) ([]dao.PlaylistModelWithUser, error) {

	return dao.GetPlaylistsByPage(page, size, curUid)
}

type UpdatePlaylistRequest struct {
	Title      string        `json:"title" binding:"required"`
	CoverURL   string        `json:"cover_url"`
	ListItem   []interface{} `json:"list_item" binding:"required"`
	PlaylistID int64         `json:"playlist_id" binding:"required"`
}

func UpdatePlaylist(req UpdatePlaylistRequest) error {
	playlist := dao.PlaylistModel{}

	playlist.Title = req.Title

	playlist.CoverURL = req.CoverURL
	listItemJSON, _ := json.Marshal(req.ListItem)
	playlist.ListItem = listItemJSON
	return dao.UpdatePlaylistByID(req.PlaylistID, playlist)
}

func DeletePlaylist(req DeletePlaylistRequest) error {
	return dao.DeletePlaylistByID(req.PlaylistID)
}

type AddPlayCountRequest struct { // 假设这是一个数组，可以根据实际情况调整类型
	PlaylistID int64 `json:"target_id" binding:"required"`
}

func AddPlayCount(req AddPlayCountRequest, curUserID int64) (bool, error) {
	play := dao.PlayCountModel{}
	return play.Add(curUserID, req.PlaylistID)
}

func DeletUserByID(userID int64) error {
	return dao.DeleteUserByID(userID)
}

func DeletePlaylistsByOwnerID(userID int64) error {
	return dao.DeletePlaylistsByOwnerID(userID)
}
