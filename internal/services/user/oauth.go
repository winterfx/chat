package user

import (
	"chat/internal/config"
	"chat/internal/repository"
	"chat/internal/repository/gen"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthProvider string

const (
	AuthProviderLocal  AuthProvider = "local"
	AuthProviderGoogle AuthProvider = "google"
)

type Auth struct {
	Provider AuthProvider
	UserRepo repository.UserRepo
	config   *oauth2.Config
}
type IdpUser struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Id    string `json:"id"`
}

func retrieveAuthConfig(p AuthProvider) (*oauth2.Config, error) {
	var clientId, secret string
	var endpoint oauth2.Endpoint
	switch p {
	case AuthProviderGoogle:
		clientId = config.Config.GoogleClientId
		secret = config.Config.GoogleClientSecret
		endpoint = google.Endpoint
	default:
		return nil, ErrInvalidAuthProvider
	}
	return &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: secret,
		RedirectURL:  fmt.Sprintf("%s/auth/google/callback", os.Getenv("FRONTEND_URL")), //"http://localhost:8080/auth/google/callback",
		Scopes:       []string{"email", "profile"},
		Endpoint:     endpoint,
	}, nil
}

func (a Auth) GetAuthConsentUrl(ctx context.Context) string {
	url := a.config.AuthCodeURL("state", oauth2.AccessTypeOffline)
	return url
}

func getUserApiByProvider(p AuthProvider) string {
	switch p {
	case AuthProviderGoogle:
		return "https://www.googleapis.com/oauth2/v2/userinfo"
	default:
		return ""
	}
}

func (a Auth) ConnectIdp(ctx context.Context, code string) (*oauth2.Token, error) {
	tok, err := a.config.Exchange(context.Background(), code)
	if err != nil {
		return nil, ErrInvalidToken.WithError(err)
	}
	return tok, nil
}

func (a Auth) RetrieveUserInfo(ctx context.Context, tok *oauth2.Token) (*IdpUser, error) {
	client := a.config.Client(context.Background(), tok)
	resp, err := client.Get(getUserApiByProvider(a.Provider))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}
	user := &IdpUser{}
	if err := json.NewDecoder(resp.Body).Decode(user); err != nil {
		return nil, err
	}
	_, err = a.UserRepo.GetUserById(ctx, user.Id)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, ErrDBOperation.WithError(err)
		} else {
			err = a.UserRepo.CreateUser(ctx, gen.CreateUserParams{
				Name:         user.Name,
				AuthProvider: a.Provider.String(),
				Email:        user.Email,
				UserID:       user.Id,
			})
			return user, nil
		}
	} else {
		err = a.UserRepo.UpdateUserById(ctx, gen.UpdateUserByIdParams{
			UserID:       user.Id,
			Email:        user.Email,
			Name:         user.Name,
			AuthProvider: a.Provider.String(),
		})
	}
	return user, nil
}

func NewOAuth(p AuthProvider, repo repository.UserRepo) (*Auth, error) {
	conf, err := retrieveAuthConfig(p)
	if err != nil {
		return nil, err
	}
	return &Auth{
		Provider: p,
		UserRepo: repo,
		config:   conf,
	}, nil
}
