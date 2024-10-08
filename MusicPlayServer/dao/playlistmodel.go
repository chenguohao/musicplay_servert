package dao

import (
	"github.com/goccy/go-json"
)

type PlaylistModel struct {
	Model
	Title      string          `json:"title"`
	PlaylistID int64           `json:"playlist_id"`
	CoverURL   string          `json:"cover_url"`
	ListItem   json.RawMessage `json:"list_item"`
	OwnerID    int64           `json:"owner_id"`
}

func (PlaylistModel) TableName() string {
	return "playlist"
}

// AddPlaylist adds a new playlist to the database
func AddNewPlaylist(playlist PlaylistModel) error {
	err := DBClient.Create(&playlist).Error
	if err != nil {
		return err
	}
	return nil
}

// DeletePlaylistsByOwnerID deletes all playlists belonging to the specified OwnerID
func DeletePlaylistsByOwnerID(ownerID int64) error {
	// 使用 DELETE 语句删除与指定 OwnerID 相关的所有播放列表
	err := DBClient.Where("owner_id = ?", ownerID).Delete(&PlaylistModel{}).Error
	if err != nil {
		return err
	}
	return nil
}

type PlaylistModelWithUser struct {
	//gorm.Model
	PlaylistModel
	Owner     SimpleUserModel `gorm:"foreignKey:OwnerID;references:Uid"` // 指定 OwnerID 对应 UserModel 的 Id
	LikeCount int             `json:"like_count"`
	PlayCount int             `json:"play_count"`
	IsLiked   bool            `json:"is_liked"`
}

func (m PlaylistModelWithUser) MarshalJSON() ([]byte, error) {
	type Alias PlaylistModelWithUser // 创建一个匿名的别名结构体，避免递归调用

	return json.Marshal(&struct {
		Alias
		CreatedAt int64 `json:"created_at"` // 将 CreatedAt 转换为 int64
	}{
		Alias:     (Alias)(m),
		CreatedAt: m.CreatedAt.Unix(), // 使用 Unix() 方法将时间转换为 Unix 时间戳
	})
}

func (PlaylistModelWithUser) TableName() string {
	return "playlist"
}
func GetPlaylistsByPage(page int, pageSize int, currentUserID int64) ([]PlaylistModelWithUser, error) {
	var playlists []PlaylistModelWithUser

	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}

	err := DBClient.Preload("Owner").
		//Find(&playlists).
		Select("playlist.*, "+
			"(SELECT COUNT(*) FROM likecount WHERE likecount.target = playlist.playlist_id and likecount.is_enable != 0 ) as like_count,"+
			"(SELECT COUNT(*) FROM playcount WHERE playcount.target = playlist.playlist_id) as play_count,"+
			"(SELECT COUNT(*) > 0 FROM likecount WHERE likecount.target = playlist.playlist_id AND likecount.sender = ? AND likecount.is_enable != 0) as is_liked", currentUserID). // 判断当前用户是否点赞 ).
		Order("playlist.created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&playlists).Error

	if err != nil {
		return nil, err
	}

	return playlists, nil
}

// GetPlaylistsByPage retrieves playlists with pagination
func GetPlaylistsByPage2(page int, pageSize int) ([]PlaylistModel, error) {
	var playlists []PlaylistModel
	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}
	err := DBClient.Offset(offset).Limit(pageSize).Find(&playlists).Error
	if err != nil {
		return nil, err
	}
	return playlists, nil
}

// UpdatePlaylistByID updates a playlist by its ID
func UpdatePlaylistByID(id int64, updatedPlaylist PlaylistModel) error {
	// 指定允许更新的字段
	err := DBClient.Model(&PlaylistModel{}).Where("playlist_id = ?", id).Updates(map[string]interface{}{
		"title":     updatedPlaylist.Title,
		"cover_url": updatedPlaylist.CoverURL,
		"list_item": updatedPlaylist.ListItem,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

// DeletePlaylistByID deletes a playlist by its ID
func DeletePlaylistByID(id int64) error {
	err := DBClient.Where("playlist_id = ?", id).Delete(&PlaylistModel{}).Error
	if err != nil {
		return err
	}
	return nil
}

// GetPlaylistCount returns the total number of playlists
func GetPlaylistCount() (int64, error) {
	var count int64
	err := DBClient.Model(&PlaylistModel{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
