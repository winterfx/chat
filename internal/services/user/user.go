package user

import (
	"chat/internal/repository"
	"chat/internal/repository/gen"
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UsersService struct {
	repo repository.UserRepo
}

func (a AuthProvider) String() string {
	return string(a)
}

func (a AuthProvider) IsLocal() bool {
	return a == AuthProviderLocal
}

func NewUsersService(repo repository.UserRepo) *UsersService {
	return &UsersService{repo: repo}
}

func (s *UsersService) userHasRegistered(ctx context.Context, email string) (bool, error) {
	_, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func encodePwd(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (s *UsersService) RegisterUser(ctx context.Context, email, password string) (string, error) {
	hasRegistered, err := s.userHasRegistered(ctx, email)
	if err != nil {
		return "", ErrDBOperation.WithError(err)
	}
	if hasRegistered {
		return "", ErrUserAlreadyRegistered
	}
	hashedPassword, err := encodePwd(password)
	if err != nil {
		return "", ErrDBOperation.WithError(err)
	}

	name := strings.Split(email, "@")[0]
	uid := uuid.New().String()
	err = s.repo.CreateUser(ctx, gen.CreateUserParams{
		UserID:       uid,
		Email:        email,
		Name:         name,
		AuthProvider: AuthProviderLocal.String(),
		PasswordHash: hashedPassword,
	})
	if err != nil {
		return "", ErrDBOperation.WithError(err)
	}
	return uid, nil
}

func (s *UsersService) Login(ctx context.Context, email, password string) (string, error) {
	r, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrUserNotFound
		}
		return "", ErrDBOperation.WithError(err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(r.PasswordHash), []byte(password))
	if err != nil {
		return "", ErrUserPasswordIncorrect
	}
	return r.UserID, nil
}

func (s *UsersService) GetUsersList(ctx context.Context, uid string) ([]*gen.User, error) {
	all, err := s.repo.GetAllUsers(ctx)
	if err != nil {
		return nil, ErrDBOperation.WithError(err)
	}
	var users []*gen.User
	for _, u := range all {
		if u.UserID == uid {
			continue
		}
		users = append(users, &gen.User{
			UserID:       u.UserID,
			Email:        u.Email,
			Name:         u.Name,
			AuthProvider: u.AuthProvider,
			PasswordHash: u.PasswordHash,
		})
	}
	return users, nil
}

func (s *UsersService) DeleteUser(ctx context.Context, uid string) error {
	err := s.repo.DeleteUserById(ctx, uid)
	if err != nil {
		return ErrDBOperation.WithError(err)
	}
	return nil
}
