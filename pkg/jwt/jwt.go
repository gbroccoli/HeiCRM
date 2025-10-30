package jwt

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	SecretKey []byte
	Issuer    string
	Audience  string
}

type FieldsClaims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

func New(secret []byte) *JWT {
	return &JWT{
		SecretKey: secret,
		Issuer:    "crm-system",
		Audience:  "zerat",
	}
}

func (j *JWT) GenerateAccessToken(email string, role int) (*string, error) {

	var obj = &FieldsClaims{
		email,
		strconv.Itoa(role),
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "user-login",
			Issuer:    "crm-system",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, obj)

	tokenString, err := token.SignedString(j.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("Ошибка при подписании токена: %v\n", err)
	}

	return &tokenString, nil
}

func (j *JWT) GenerateRefreshToken() string {
	return ""
}

func (j *JWT) Verify(token string) (*FieldsClaims, error) {

	tokenParser, err := jwt.ParseWithClaims(token, &FieldsClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.SecretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("Ошибка при парсинге токена: %v\n", err)
	}

	if claims, ok := tokenParser.Claims.(*FieldsClaims); ok && tokenParser.Valid {
		return claims, nil
	}

	return nil, errors.New("no valid token claims")
}
