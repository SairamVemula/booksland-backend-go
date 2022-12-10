package middlewares

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/SairamVemula/booksland-backend-go/pkg/models"
	"github.com/SairamVemula/booksland-backend-go/pkg/services"
	"github.com/SairamVemula/booksland-backend-go/pkg/utils"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//authenticatedHandler is a handler function that also requires a user
type AuthenticatedHandler func(http.ResponseWriter, *http.Request, *models.User)

type EnsureAuth struct {
	handler    AuthenticatedHandler
	roles      []string
	middleware *Middleware
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func (ea *EnsureAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var user models.User
	// log.Println(r.URL)

	tokenString := r.Header.Get("Authorization")
	// log.Println(tokenString)
	if len(tokenString) == 0 && len(ea.roles) != 0 {
		utils.ResponseError(&w, &utils.RestError{Code: http.StatusUnauthorized, Message: "Missing Authorization Header"})
		return
	}
	if len(tokenString) != 0 {
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		user_id, role, err := ea.middleware.ValidateAccessToken(tokenString)
		if err != nil {
			utils.ResponseError(&w, &utils.RestError{Code: http.StatusUnauthorized, Message: "Error verifying JWT token: " + err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
		defer cancel()
		id, e := primitive.ObjectIDFromHex(user_id)
		if e != nil {
			utils.ResponseError(&w, &utils.RestError{Code: http.StatusUnauthorized, Message: "Invalid Authorization Header"})
			return
		}
		err = models.UsersCollection.FindOne(ctx, bson.M{"_id": id, "type": role}).Decode(&user)
		if err != nil {
			utils.ResponseError(&w, &utils.RestError{Code: http.StatusUnauthorized, Message: "Invalid Authorization Header"})
			return
		}
	}
	if user.Phone == "" {
		ea.handler(w, r, nil)
		return
	}
	if len(ea.roles) != 0 && !contains(ea.roles, user.Type) {
		utils.ResponseError(&w, &utils.RestError{Code: http.StatusForbidden, Message: "Forbbiden Access"})
		return
	}
	ea.handler(w, r, &user)
}

func (m *Middleware) AuthWithRoles(roles []string, handlerToWrap AuthenticatedHandler) *EnsureAuth {
	return &EnsureAuth{handlerToWrap, roles, m}
}

func (m *Middleware) ValidateAccessToken(signedToken string) (string, string, error) {

	token, err := jwt.ParseWithClaims(signedToken, &services.AccessTokenCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			m.logger.Error("Unexpected signing method in auth token")
			return nil, errors.New("Unexpected signing method in auth token")
		}
		verifyBytes, err := ioutil.ReadFile(m.configs.AccessTokenPublicKeyPath)
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
		m.logger.Error("unable to parse claims", "error", err)
		return "", "", err
	}

	claims, ok := token.Claims.(*services.AccessTokenCustomClaims)
	if !ok || !token.Valid || claims.UserID == "" || claims.KeyType != "access" {
		return "", "", errors.New("invalid token: authentication failed")
	}
	return claims.UserID, claims.UserType, nil
}
func (m *Middleware) ValidateRefreshToken(tokenString string) (string, error) {

	token, err := jwt.ParseWithClaims(tokenString, &services.RefreshTokenCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
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
		return "", err
	}

	claims, ok := token.Claims.(*services.RefreshTokenCustomClaims)
	m.logger.Debug("ok", ok)
	if !ok || !token.Valid || claims.UserID == "" || claims.KeyType != "refresh" {
		m.logger.Debug("could not extract claims from token")
		return "", errors.New("invalid token: authentication failed")
	}
	return claims.UserID, nil
}
