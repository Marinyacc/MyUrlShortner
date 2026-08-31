package router

import (
	"html/template"
	"my_url_shortner/global"
	"my_url_shortner/handler"
	"my_url_shortner/middleware"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func templateFuncMap() template.FuncMap {
	return template.FuncMap{
		"portalBaseURL": func() string {
			return trimBaseURL(global.Conf.App.UrlPrefix)
		},
		"adminBaseURL": func() string {
			return trimBaseURL(global.Conf.App.AdminUrlPrefix)
		},
	}
}

func trimBaseURL(raw string) string {
	return strings.TrimRight(strings.TrimSpace(raw), "/")
}

//localhost:9092/
func NewAdminRouter() *gin.Engine {
	r := gin.Default()
	r.SetFuncMap(templateFuncMap())
	r.LoadHTMLGlob("templates/*")
	// r.Use(middleware.LogMiddleware())

	r.GET("/login", handler.LoginPage)               //登陆界面
	r.POST("/login", handler.Login)                  //尝试登陆
	r.GET("/captcha/:id", handler.ServeCaptchaImage) //获取验证码图片
	r.GET("/captcha", handler.RequestCaptchaImage)   //获取验证码ID

	admin := r.Group("/admin")
	admin.Use(middleware.JWTAuthMiddleware())
	admin.POST("/logout", handler.Logout)                  //登出
	admin.GET("/dashboard", handler.DashboardPage)         //获取仪表盘界面
	admin.GET("/urls", handler.UrlsPage)                   //获取短链接列表
	admin.GET("/stats", handler.StatsPage)                 //获取短链接数据统计列表
	admin.GET("/access_logs", handler.AccessLogsPage)      //获取访问日志列表
	admin.POST("/urls/generate", handler.GenerateShortUrl) //新建短链接
	admin.POST("/urls/delete", handler.DeleteShortUrl)     //删除短链接
	// admin.POST("/access_logs_export", handler.AccessLogsExport) //导出访问日志

	api := r.Group("/api")
	// api.Use(middleware.JWTAuthMiddleware())
	api.GET("/account", handler.APIAccountInfo)        //获取所有用户信息
	api.POST("/account", handler.APINewUser)           //新建一个用户
	api.POST("/account/update", handler.APIUpdateUser) //更新一个用户的密码
	api.POST("/account/delete", handler.APIDeleteUser) //删除一个用户

	api.GET("/urls/:url", handler.APIUrlInfo)        //获取指定的短链接的信息
	api.GET("/urls", handler.APIUrlsInfo)            //获取所有短链接信息的列表
	api.GET("/urls/stats/:url", handler.APIUrlStats) //获取指定的短链接的统计信息
	api.GET("/urls/stats", handler.APIUrlsStats)     // 获取所有短链接的统计信息

	api.POST("/url", handler.APIGenShortUrl) //生成短链接

	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "error.gohtml", gin.H{
			"error": "不存在的页面",
		})
	})
	return r
}

//localhost:9091/
func NewVistorRouter() *gin.Engine {
	r := gin.Default()
	r.SetFuncMap(templateFuncMap())
	r.LoadHTMLGlob("templates/*")
	r.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	r.GET("/:url", handler.RedirectLongUrl) //通过shorturl跳转至目标页面
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "error.gohtml", gin.H{
			"error": "不存在的页面",
		})
	})
	return r
}
