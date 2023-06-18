package utils

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"math/rand"
	"site/types"
	"time"
)

const (
	salt            = "hfasdgfhasdgfshjdggf"
	signingKey      = "%*FG67G%f786^G%&()(&J*H)(_I*K{76534d5D"
	tokenTTL        = 15 * time.Second
	refreshTokenTTL = 5 * 24 * time.Hour
)

type userClaims struct {
	jwt.StandardClaims
	User types.User `json:"user"`
}

func GeneratePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func NewRefreshToken() (string, int64, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", 0, err
	}

	return fmt.Sprintf("%x", b), time.Now().Add(refreshTokenTTL).Unix(), nil
}

func NewSessionUserToken(user *types.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &userClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		*user,
	})

	return token.SignedString([]byte(signingKey))
}

func ParseUserToken(accessToken string) (*types.User, error) {
	token, err := jwt.ParseWithClaims(accessToken, &userClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*userClaims)
	if !ok {
		return nil, errors.New("token claims are not of type *userClaims")
	}

	return &claims.User, nil
}
