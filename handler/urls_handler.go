package handler

import (
	"my_url_shortner/global"
	"my_url_shortner/storage"
	"my_url_shortner/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// RedirectLongUrl 根据短链接跳转至长链接网址
func RedirectLongUrl(c *gin.Context) {
	shortUrl := c.Param("url")
	if utils.EmptyString(shortUrl) {
		c.HTML(http.StatusBadRequest, "error.gohtml", gin.H{
			"error": "短链接为空",
		})
		return
	}

	urlKey := global.ShortUrl + shortUrl
	res, err := storage.RedisGetString(c.Request.Context(), urlKey)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.gohtml", gin.H{
			"error": "Redis 获取键值对错误",
		})
		return
	} else if res == "NULL" {
		c.HTML(http.StatusNotFound, "error.gohtml", gin.H{
			"error": "长链接地址不存在",
		})
		return
	}

	//  Redis 没查到
	if utils.EmptyString(res) {
		lockKey := global.ShortUrl_Lock + shortUrl
		lock := storage.NewLock(c.Request.Context(), lockKey, 2*time.Second)

		//抢到锁
		if lock.Lock(c.Request.Context()) {
			defer lock.UnLock(c.Request.Context())

			//check
			longUrl, err := storage.RedisGetString(c.Request.Context(), urlKey)
			if err == nil && longUrl != "NULL" && !utils.EmptyString(longUrl) {
				RedirectSuccess(c, shortUrl, longUrl)
				return
			}

			//数据库中查询
			url, err := storage.GetUrlInfo(c.Request.Context(), shortUrl)
			if err != nil {
				storage.RedisSet30m(c.Request.Context(), urlKey, "NULL")
				c.HTML(http.StatusNotFound, "error.gohtml", gin.H{
					"error": "长连接未在数据库中找到",
				})
				return
			}

			//回写Redis
			if err := storage.RedisSet4Ever(c.Request.Context(), urlKey, url.LongUrl); err != nil {
				c.HTML(http.StatusInternalServerError, "error.gohtml", gin.H{
					"error": "Redis回写错误",
				})
				return
			}
			RedirectSuccess(c, shortUrl, url.LongUrl)
			return
		}
		//没抢到锁
		time.Sleep(50 * time.Millisecond)
		retryRes, err := storage.RedisGetString(c.Request.Context(), urlKey)
		if err == nil && !utils.EmptyString(retryRes) && retryRes != "NULL" {
			c.Redirect(http.StatusFound, retryRes)
			return
		}
		// 依然没有则返回 404 或重试
		c.HTML(http.StatusNotFound, "error.gohtml", gin.H{
			"error": "未找到长连接,可重新尝试",
		})
	} else {
		//Redis查到
		RedirectSuccess(c, shortUrl, res)
	}
}

func RedirectSuccess(c *gin.Context, shortUrl, longUrl string) {
	c.Redirect(http.StatusFound, longUrl)

	go storage.Ip_Count_Incr(shortUrl, c.ClientIP())
	go storage.ProduceMessage(global.StreamName, shortUrl, c.ClientIP())
}
