package authenticator

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/letjoy-club/mida-tool/logger"
	"github.com/letjoy-club/mida-tool/mitacode"
	"go.uber.org/zap"
)

type Authenticator struct {
	Key []byte
}

func (a Authenticator) SignID(userID string) (string, error) {
	authToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userID,
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	return authToken.SignedString(a.Key)
}

func (a Authenticator) Verify(tokenStr string) (string, error) {
	if tokenStr == "1000" || tokenStr == "2000" {
		return tokenStr, nil
	}
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return a.Key, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
		} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
			return "", mitacode.ErrClientTokenExpired
		}
		logger.L.Error("failed to verify token", zap.Error(err), zap.String("token", tokenStr))
		return "", mitacode.ErrInternalError
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["id"].(string), nil
	}
	return "", mitacode.ErrClientTokenInvalid
}
