package repository

import (
	"chat/internal/repository/gen"
	"context"
	"database/sql"
)

type UserRepo interface {
	CreateUser(ctx context.Context, arg gen.CreateUserParams) error
	GetUserById(ctx context.Context, userID string) (gen.GetUserByIdRow, error)
	GetUserByEmail(ctx context.Context, email string) (gen.GetUserByEmailRow, error)
	UpdateUserById(ctx context.Context, arg gen.UpdateUserByIdParams) error
	UpdateUserByEmail(ctx context.Context, arg gen.UpdateUserByEmailParams) error
	GetAllUsers(ctx context.Context) ([]gen.GetAllUsersRow, error)
	DeleteUserById(ctx context.Context, userID string) error
}

type userRepo struct {
	query *gen.Queries
}

func NewUserRepo(db *sql.DB) UserRepo {
	return &userRepo{query: gen.New(db)}
}

func (u *userRepo) GetAllUsers(ctx context.Context) ([]gen.GetAllUsersRow, error) {
	return u.query.GetAllUsers(ctx)
}

func (u *userRepo) CreateUser(ctx context.Context, arg gen.CreateUserParams) error {
	return u.query.CreateUser(ctx, arg)
}

func (u *userRepo) GetUserById(ctx context.Context, userID string) (gen.GetUserByIdRow, error) {
	return u.query.GetUserById(ctx, userID)
}

func (u *userRepo) GetUserByEmail(ctx context.Context, email string) (gen.GetUserByEmailRow, error) {
	return u.query.GetUserByEmail(ctx, email)
}

func (u *userRepo) UpdateUserById(ctx context.Context, arg gen.UpdateUserByIdParams) error {
	return u.query.UpdateUserById(ctx, arg)
}

func (u *userRepo) UpdateUserByEmail(ctx context.Context, arg gen.UpdateUserByEmailParams) error {
	return u.query.UpdateUserByEmail(ctx, arg)
}

func (u *userRepo) DeleteUserById(ctx context.Context, userID string) error {
	return u.query.DeleteUserById(ctx, userID)
}
