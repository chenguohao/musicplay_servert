package controller

import (
	. "MusicPlayServer/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func TestApi(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Authentication successful",
		"code":    0,
	})
}

func GetPlaylist(c *gin.Context) {

	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))

	if size == 0 {
		size = 1
	}

	userIDHeader := c.GetHeader("X-User-ID")

	// 如果 userIDHeader 为空，返回错误
	if userIDHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in header"})
		return
	}

	curUID, _ := strconv.Atoi(userIDHeader)
	if playlist, err := GetPlayList(page, size, int64(curUID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": 1001})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Playlist created successfully",
			"data":    playlist,
			"code":    0,
		})
	}

}

func CreatePlayList(c *gin.Context) {
	var req CreatePlaylistRequest

	// 绑定传入的 JSON 数据到 PlaylistModel 结构体
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 生成唯一的 PlaylistID
	if err := CreatePlaylist(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": 1001})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"message": "Playlist created successfully",
	})
}

func AuthWithApple(c *gin.Context) {
	var req AppleLoginRequest

	// 绑定 JSON 参数并验证
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": 1001})
		return
	}

	result := AppleSign(req.AuthorizationCode, req.IDToken)

	if !result {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Authentication fail",
			"code":    1001,
		})
		return
	}
	userInfo := RegisterOrLogin(req)
	c.JSON(http.StatusOK, gin.H{
		"message": "Authentication successful",
		"data":    userInfo,
		"code":    0,
	})
}

func ReqestUpdateProfile(c *gin.Context) {
	var req ProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": 1001})
		return
	}

	if err := UpdateProfile(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": 1001})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Authentication successful",
		"data":    "",
		"code":    0,
	})
}

func RequestUpdatePlaylist(c *gin.Context) {
	var req UpdatePlaylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": 1001})
		return
	}

	if err := UpdatePlaylist(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": 1001})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "update successful",
		"data":    "",
		"code":    0,
	})
}

func RequestDeletePlaylist(c *gin.Context) {
	var req DeletePlaylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": 1001})
		return
	}

	if err := DeletePlaylist(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": 1001})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Authentication successful",
		"data":    "",
		"code":    0,
	})
}

func RequestLike(c *gin.Context) {
	var req LikeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": 1001})
		return
	}

	userIDHeader := c.GetHeader("X-User-ID")

	// 如果 userIDHeader 为空，返回错误
	if userIDHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in header"})
		return
	}

	curUID, _ := strconv.Atoi(userIDHeader)

	if err := DoLike(req, int64(curUID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": 1001})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Authentication successful",
		"data":    "",
		"code":    0,
	})
}

func RequestAddPlayCount(c *gin.Context) {
	var req AddPlayCountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": 1001})
		return
	}

	userIDHeader := c.GetHeader("X-User-ID")

	// 如果 userIDHeader 为空，返回错误
	if userIDHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in header"})
		return
	}

	curUID, _ := strconv.Atoi(userIDHeader)

	if isAdd, err := AddPlayCount(req, int64(curUID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": 1001})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Authentication successful",
			"data":    map[string]interface{}{"isAdd": isAdd},
			"code":    0,
		})
		return
	}

}
