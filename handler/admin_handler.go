package handler

import (
	"my_url_shortner/global"
	"my_url_shortner/middleware"
	"my_url_shortner/model"
	"my_url_shortner/service"
	"my_url_shortner/storage"
	"my_url_shortner/utils"
	"net/http"
	"strconv"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

// LoginPage 登陆界面
func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.gohtml", gin.H{})
}

// Login 登陆请求
func Login(c *gin.Context) {
	account := c.PostForm("account")
	password := c.PostForm("password")
	captchaText := c.PostForm("captcha-text")
	captchaID := c.PostForm("captcha-id")

	//用户账号密码合法性验证
	if utils.EmptyString(account) || utils.EmptyString(password) {
		c.HTML(http.StatusOK, "login.gohtml", gin.H{
			"error": "用户名或密码不能为空",
		})
		return
	}

	//6位验证码合法性验证
	if utils.EmptyString(captchaText) || utils.EmptyString(captchaID) || len(captchaText) != 6 {
		c.HTML(http.StatusOK, "login.gohtml", gin.H{
			"error": "验证码格式不正确",
		})
		return
	}
	if !captcha.VerifyString(captchaID, captchaText) {
		c.HTML(http.StatusOK, "login.gohtml", gin.H{
			"error": "验证码不正确",
		})
		return
	}

	if service.Login(c.Request.Context(), account, password) {
		jwtToken, err := middleware.GenerateToken(account)
		if err == nil {
			c.SetCookie("token", jwtToken, global.CookieTTL, "/", "", false, true)
		}
		c.Redirect(http.StatusFound, "/admin/dashboard")
	} else {
		c.HTML(http.StatusOK, "login.gohtml", gin.H{
			"error": "用户名或密码不正确",
		})
	}
}

// ServeCaptchaImage 生成验证码
func ServeCaptchaImage(c *gin.Context) {
	captcha.Server(200, 45).ServeHTTP(c.Writer, c.Request)
}

// RequestCaptchaImage 生成验证码图片
func RequestCaptchaImage(c *gin.Context) {
	imageID := captcha.New()
	c.JSON(http.StatusOK, model.ResultJsonSuccessWithData(imageID))
}

// Logout 登出
func Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/login")
}

// DashboardPage 仪表盘界面
func DashboardPage(c *gin.Context) {
	query := `SELECT * FROM public.stats_sum`

	var res []model.StatsSum
	err := storage.DbSelect(c.Request.Context(), query, &res)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.gohtml", gin.H{
			"error": err.Error(),
		})
		return
	}

	statsMap := make(map[string]int, len(res))
	for _, stats := range res {
		statsMap[stats.Key] = stats.Value
	}

	c.HTML(http.StatusOK, "dashboard.gohtml", gin.H{
		"total_pv":            statsMap["total_pv"],
		"total_uv":            statsMap["total_uv"],
		"today_count":         statsMap["today_count"],
		"yesterday_count":     statsMap["yesterday_count"],
		"last_7_days_count":   statsMap["last_7_days_count"],
		"monthly_count":       statsMap["monthly_count"],
		"d_today_count":       statsMap["d_today_count"],
		"d_yesterday_count":   statsMap["d_yesterday_count"],
		"d_last_7_days_count": statsMap["d_last_7_days_count"],
		"d_monthly_count":     statsMap["d_monthly_count"],
	})
}

// UrlsPage 短链接列表界面
func UrlsPage(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	totalCount, _ := storage.TotalUrls(c.Request.Context())

	if page < 1 {
		page = global.Page
	}
	if pageSize > 100 || pageSize < 1 {
		pageSize = global.PageSize
	}
	totalPage := (totalCount + pageSize - 1) / pageSize
	if totalPage < 1 {
		totalPage = 1
	}

	shortUrls, err := storage.GetUrlsInfoByPage(c.Request.Context(), page, pageSize)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.gohtml", gin.H{
			"error": err.Error(),
		})
		return
	}
	c.HTML(http.StatusOK, "urls.gohtml", gin.H{
		"short_urls": shortUrls,
		"page":       page,
		"pageSize":   pageSize,
		"totalPage":  totalPage,
	})
}

// StatsPage 短链接统计列表界面
func StatsPage(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	totalCount, _ := storage.TotalStats(c.Request.Context())

	if page < 1 {
		page = global.Page
	}
	if pageSize > 100 || pageSize < 1 {
		pageSize = global.PageSize
	}
	totalPage := (totalCount + pageSize - 1) / pageSize
	if totalPage < 1 {
		totalPage = 1
	}

	shortUrlsStats, err := storage.GetUrlsStatsByPage(c.Request.Context(), page, pageSize)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.gohtml", gin.H{
			"error": err.Error(),
		})
		return
	}
	c.HTML(http.StatusOK, "urls_stats.gohtml", gin.H{
		"short_url_stats": shortUrlsStats,
		"page":            page,
		"pageSize":        pageSize,
		"totalPage":       totalPage,
	})
}

// AccessLogsPage 访问日志列表界面
func AccessLogsPage(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	totalCount, _ := storage.TotalAccessLogs(c.Request.Context())

	if page < 1 {
		page = global.Page
	}
	if pageSize > 100 || pageSize < 1 {
		pageSize = global.PageSize
	}
	totalPage := (totalCount + pageSize - 1) / pageSize
	if totalPage < 1 {
		totalPage = 1
	}

	accessLogs, err := storage.GetAccessLogsByPage(c.Request.Context(), page, pageSize)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.gohtml", gin.H{
			"error": err.Error(),
		})
		return
	}
	c.HTML(http.StatusOK, "access_logs.gohtml", gin.H{
		"logs":      accessLogs,
		"page":      page,
		"pageSize":  pageSize,
		"totalPage": totalPage,
	})
}

// GenerateShortUrl 创建新的短链接
func GenerateShortUrl(c *gin.Context) {
	longUrl := c.PostForm("long_url")
	comment := c.PostForm("comment")

	if utils.EmptyString(longUrl) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "longUrl cannot be empty",
		})
		return
	}

	shortUrl, err := service.GenerateShortUrl(c.Request.Context(), longUrl, comment)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, shortUrl)
}

// DeleteShortUrl 删除短链接
func DeleteShortUrl(c *gin.Context) {
	shortUrl := c.PostForm("short_url")

	query1 := `DELETE FROM public.short_urls WHERE short_url = $1`
	query2 := `DELETE FROM public.stats WHERE short_url = $1`
	query3 := `DELETE FROM public.access_logs WHERE short_url = $1`
	query4 := `DELETE FROM public.daily_stats WHERE short_url = $1`

	err := storage.DbExecTx(c.Request.Context(),
		storage.Query{Query: query1, Args: []any{shortUrl}},
		storage.Query{Query: query2, Args: []any{shortUrl}},
		storage.Query{Query: query3, Args: []any{shortUrl}},
		storage.Query{Query: query4, Args: []any{shortUrl}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	storage.RedisDelete(c.Request.Context(), global.ShortUrl+shortUrl)
	//todo

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
