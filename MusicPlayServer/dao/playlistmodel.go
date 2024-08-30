package dao

import "github.com/goccy/go-json"

type PlaylistModel struct {
	ID         int64           `json:"id"`
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
func AddPlaylist(playlist PlaylistModel) error {
	err := DBClient.Create(&playlist).Error
	if err != nil {
		return err
	}
	return nil
}

// GetPlaylistsByPage retrieves playlists with pagination
func GetPlaylistsByPage(page int, pageSize int) ([]PlaylistModel, error) {
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
	err := DBClient.Model(&PlaylistModel{}).Where("id = ?", id).Updates(updatedPlaylist).Error
	if err != nil {
		return err
	}
	return nil
}

// DeletePlaylistByID deletes a playlist by its ID
func DeletePlaylistByID(id int64) error {
	err := DBClient.Where("id = ?", id).Delete(&PlaylistModel{}).Error
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
