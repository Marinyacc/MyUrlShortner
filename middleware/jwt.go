package middleware

import (
	"errors"
	"fmt"
	"log"
	"my_url_shortner/global"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	Account string `json:"account"`
	jwt.RegisteredClaims
}

// GenerateToken 根据账户生成 JWT 字符串
func GenerateToken(account string) (string, error) {
	claims := CustomClaims{
		Account: account,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "my_url_shortner",
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return jwtToken.SignedString([]byte(os.Getenv("SECRET")))
}

// ParseToken 返回token字符串的payload
func ParseToken(token string) (*CustomClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("未知的签名算法:%v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := jwtToken.Claims.(*CustomClaims); ok && jwtToken.Valid {
		return claims, nil
	}
	return nil, errors.New("非法的Token")
}

// JWT 中间件,验证权限
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		claims, err := ParseToken(token)
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		timeRemain := time.Until(claims.ExpiresAt.Time)
		if timeRemain < 5*time.Minute {
			newToken, err := GenerateToken(claims.Account)
			if err == nil {
				c.SetCookie("token", newToken, global.CookieTTL, "/", "", false, true)
			} else {
				log.Println("JWT token 滑动窗口续期发生错误", err)
			}
		}

		c.Set("account", claims.Account)
		c.Next()
	}
}
