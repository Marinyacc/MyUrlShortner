package storage

import (
	"context"
	"my_url_shortner/model"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

// FindAllUsers 获取所有用户
func FindAllUsers(ctx context.Context) ([]model.User, error) {
	query := `SELECT * FROM public.users`
	var found []model.User
	return found, DbSelect(ctx, query, &found)
}

// NewUser 新建用户
func NewUser(ctx context.Context, account string, password string) error {
	query := `INSERT INTO public.users (account,"password") VALUES (@account,@password)`
	pwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return DbNamedExec(ctx, query, pgx.NamedArgs{
		"account":  account,
		"password": string(pwd),
	})
}

// UpdateUser 更新用户
func UpdateUser(ctx context.Context, account, password string) error {
	query := `UPDATE public.users SET "password" = @password WHERE account = @account`
	return DbNamedExec(ctx, query, pgx.NamedArgs{
		"account":  account,
		"password": password,
	})
}

// DeleteUser 删除用户
func DeleteUser(ctx context.Context, user model.User) error {
	query := `DELETE FROM public.users WHERE id = @id`
	return DbNamedExec(ctx, query, pgx.NamedArgs{
		"id": user.ID,
	})
}

// FindUserByAccount 根据账户查找用户
func FindUserByAccount(ctx context.Context, account string) (model.User, error) {
	query := `SELECT * FROM public.users u WHERE lower(u.account)= $1`
	var user model.User
	return user, DbGet(ctx, query, &user, account)
}
