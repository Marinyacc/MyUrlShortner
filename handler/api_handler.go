package handler

import (
	"log"
	"my_url_shortner/service"
	"my_url_shortner/storage"
	"my_url_shortner/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type NewUserRequest struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserRequest struct {
	Account     string `json:"account" binding:"required"`
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type DeleteUserRequest struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type NewUrlRequest struct {
	LongUrl string `json:"long_url" binding:"required"`
	Comment string `json:"comment" binding:"required"`
}

// APINewUser 创建新用户
func APINewUser(c *gin.Context) {
	var req NewUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求格式错误",
		})
		return
	}
	if err := storage.NewUser(c.Request.Context(), req.Account, req.Password); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "数据库发生错误",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "用户创建成功",
	})
}

// APIAccountInfo 获取所有用户信息
func APIAccountInfo(c *gin.Context) {
	users, err := storage.FindAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, users)
}

// APIUpdateUser 更新用户密码
func APIUpdateUser(c *gin.Context) {
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	user, err := storage.FindUserByAccount(c.Request.Context(), req.Account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "原始密码不正确",
		})
		return
	}

	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := storage.UpdateUser(c.Request.Context(), req.Account, string(newPasswordHash)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "更新用户信息成功",
	})
}

// APIDeleteUser 删除用户
func APIDeleteUser(c *gin.Context) {
	var req DeleteUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	user, err := storage.FindUserByAccount(c.Request.Context(), req.Account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := storage.DeleteUser(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "删除用户成功",
	})
}

// APIUrlInfo 获取指定的短链接的信息
func APIUrlInfo(c *gin.Context) {
	shortUrl := c.Param("url")

	UrlInfo, err := storage.GetUrlInfo(c.Request.Context(), shortUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, UrlInfo)
}

// APIUrlsInfo 获取所有短链接统计信息
func APIUrlsInfo(c *gin.Context) {
	UrlsInfo, err := storage.GetUrlsInfo(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, UrlsInfo)
}

// APIUrlState 获取指定Url的统计信息
func APIUrlStats(c *gin.Context) {
	shorUrl := c.Param("url")
	Urlstate, err := storage.GetUrlStats(c.Request.Context(), shorUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, Urlstate)
}

// APIUrlsState 获取所有Url的统计信息
func APIUrlsStats(c *gin.Context) {
	Urlstate, err := storage.GetUrlsStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, Urlstate)
}

// APIGenShortUrl 生成一个新的短链接
func APIGenShortUrl(c *gin.Context) {
	var req NewUrlRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if utils.EmptyString(req.LongUrl) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "原始长连接为空",
		})
		return
	}

	shortUrl, err := service.GenerateShortUrl(c.Request.Context(), req.LongUrl, req.Comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "短链接创建成功",
		"short_url": shortUrl,
	})
}
