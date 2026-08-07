package service

import (
	"context"
	"my_url_shortner/storage"

	"golang.org/x/crypto/bcrypt"
)

// Login 登陆
func Login(ctx context.Context, account string, password string) bool {
	User, err := storage.FindUserByAccount(ctx, account)
	if err != nil {
		return false
	}
	err = bcrypt.CompareHashAndPassword([]byte(User.Password), []byte(password))
	return err == nil
}
