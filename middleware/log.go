package middleware

// func LogMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		temp := time.Now()
// 		log.Println("用户IP地址:", c.ClientIP())
// 		c.Next()
// 		log.Println("调用时间:",time.Since(temp))
// 	}
// }
