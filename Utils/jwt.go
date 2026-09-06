package Utils

import (
	"MyGO-IM/Conf"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

var MySecret = []byte(Conf.Conf.J.Secret)

func CreateJWT(userName string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "RippleHeart",
		Subject:   "OK",
		Audience:  []string{userName},
		ID:        uuid.NewString(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(Conf.Conf.J.TTL))),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})
	tokenStr, err := token.SignedString(MySecret)
	if err != nil {
		return "", nil
	}
	return tokenStr, nil
}
func VerifyJWT(token string) (string, bool) {
	tokenParse, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if token.Method == jwt.SigningMethodHS256 {
			return MySecret, nil
		}
		return nil, errors.New("error parse method")
	}, jwt.WithValidMethods([]string{"HS256"}))
	if tokenParse == nil || !tokenParse.Valid || err != nil {
		return "", false
	}
	audience, err := tokenParse.Claims.GetAudience()
	if err != nil {
		return "", false
	}
	return audience[0], true
}
