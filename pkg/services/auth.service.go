package services

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/SairamVemula/booksland-backend-go/pkg/models"
	"github.com/SairamVemula/booksland-backend-go/pkg/utils"
	"github.com/golang-jwt/jwt"
	"github.com/hashicorp/go-hclog"
)

var as *AuthService

type AuthService struct {
	us        *UserService
	logger    hclog.Logger
	configs   *utils.Configurations
	validator *models.Validation
}

func NewAuthService(logger hclog.Logger, configs *utils.Configurations, validator *models.Validation) *AuthService {
	if as != nil {
		return as
	}
	us := NewUserService(logger, configs, validator)
	return &AuthService{us, logger, configs, validator}
}

func (as *AuthService) Register(ctx context.Context, user *models.User) (*models.User, *utils.RestError) {
	return as.us.Create(ctx, user)
}

func (as *AuthService) Login(ctx context.Context, loginUser *models.LoginUser) (*models.User, *utils.RestError) {
	user, err := as.us.FindUsernameAndPassword(ctx, loginUser.Username, loginUser.Password)
	if err != nil {
		return nil, err
	}
	token, terr := as.GenerateAccessToken(ctx, user.ID.Hex(), user.Type)
	if terr != nil {
		return nil, &utils.RestError{
			Message: "Error on Token generation",
			Code:    http.StatusInternalServerError,
			Error:   terr.Error(),
		}
	}
	refresh, terr := as.GenerateRefreshToken(ctx, user.ID.Hex())
	if terr != nil {
		return nil, &utils.RestError{
			Message: "Error on Token generation",
			Code:    http.StatusInternalServerError,
			Error:   terr.Error(),
		}
	}
	user.Token = token
	user.TokenExpiry = (time.Now().UnixMilli() + int64(as.configs.JwtExpiration*60000)) - (15 * 60 * 1000)
	user.RefreshToken = refresh
	user.RefreshTokenExpiry = (time.Now().UnixMilli() + int64(as.configs.RefreshJwtExpiration*60000))
	user.Password = ""
	return user, nil
}

func (as *AuthService) StoreRefreshToken(ctx context.Context, user_id string, token string) (*models.User, *utils.RestError) {
	return as.us.UpdateById(ctx, user_id, &models.UpdateUser{RefreshToken: token})
}

// RefreshTokenCustomClaims specifies the claims for refresh token
type RefreshTokenCustomClaims struct {
	UserID  string `json:"user_id,omitempty"`
	KeyType string `json:"key_type,omitempty"`
	jwt.StandardClaims
}

// AccessTokenCustomClaims specifies the claims for access token
type AccessTokenCustomClaims struct {
	UserID   string `json:"user_id,omitempty"`
	KeyType  string `json:"key_type,omitempty"`
	UserType string `json:"user_type,omitempty"`
	jwt.StandardClaims
}

var jwtKey = []byte("supersecretkey")
var jwtRefreshKey = []byte("supersecretrefreshkey")

func (as *AuthService) GenerateAccessToken(ctx context.Context, user_id string, user_type string) (tokenString string, err error) {

	claims := AccessTokenCustomClaims{
		user_id,
		"access",
		user_type,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(as.configs.JwtExpiration)).Unix(),
			Issuer:    "booksland.auth.service",
		},
	}

	signBytes, err := ioutil.ReadFile(as.configs.AccessTokenPrivateKeyPath)
	if err != nil {
		as.logger.Error("unable to read private key", "error", err)
		return "", errors.New("could not generate refresh token. please try again later")
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		as.logger.Error("unable to parse private key", "error", err)
		return "", errors.New("could not generate refresh token. please try again later")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(signKey)
}

func (as *AuthService) GenerateRefreshToken(ctx context.Context, user_id string) (tokenString string, err error) {

	claims := RefreshTokenCustomClaims{
		user_id,
		"refresh",
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(as.configs.RefreshJwtExpiration)).Unix(),
			Issuer:    "booksland.auth.service",
		},
	}

	signBytes, err := ioutil.ReadFile(as.configs.RefreshTokenPrivateKeyPath)
	if err != nil {
		as.logger.Error("unable to read private key", "error", err)
		return "", errors.New("could not generate refresh token. please try again later")
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		as.logger.Error("unable to parse private key", "error", err)
		return "", errors.New("could not generate refresh token. please try again later")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenString, err = token.SignedString(signKey)

	as.StoreRefreshToken(ctx, user_id, tokenString)

	return
}

func (m *AuthService) ValidateRefreshToken(ctx context.Context, tokenString string) (string, string, error) {

	token, err := jwt.ParseWithClaims(tokenString, &RefreshTokenCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			m.logger.Error("Unexpected signing method in auth token")
			return nil, errors.New("Unexpected signing method in auth token")
		}
		verifyBytes, err := ioutil.ReadFile(m.configs.RefreshTokenPublicKeyPath)
		if err != nil {
			m.logger.Error("unable to read public key", "error", err)
			return nil, err
		}

		verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
		if err != nil {
			m.logger.Error("unable to parse public key", "error", err)
			return nil, err
		}

		return verifyKey, nil
	})

	if err != nil {
		// m.logger.Error("unable to parse claims", "error", err)
		return "", "", err
	}

	claims, ok := token.Claims.(*RefreshTokenCustomClaims)
	m.logger.Debug("ok", ok)
	if !ok || !token.Valid || claims.UserID == "" || claims.KeyType != "refresh" {
		m.logger.Debug("could not extract claims from token")
		return "", "", errors.New("invalid token: authentication failed")
	}

	user, errf := m.us.FindById(ctx, claims.UserID)
	if errf != nil {
		return "", "", errors.New(errf.Message)
	}
	if user.RefreshToken != tokenString {
		return "", "", errors.New("invalid token: authentication failed")
	}

	return claims.UserID, user.Type, nil
}

func (as *AuthService) Logout(ctx context.Context, user *models.User) {
	as.us.UpdateById(ctx, user.ID.Hex(), &models.UpdateUser{RefreshToken: "  "})
}
