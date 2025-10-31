package jwt

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
	//ErrInvalidSignature = errors.New("invalid token signature")
)

type JWT struct {
	SecretKey []byte
	Issuer    string
	Audience  string
}

type FieldsClaims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	Type  string `json:"type"`
	jwt.RegisteredClaims
}

func New(secret []byte) *JWT {
	return &JWT{
		SecretKey: secret,
		Issuer:    "crm-system",
		Audience:  "zerat",
	}
}

func (j *JWT) GenerateAccessToken(email string, role int, tokenType string) (string, error) {

	if email == "" {
		return "", fmt.Errorf("email cannot be empty")
	}

	if tokenType == "" {
		tokenType = "access" // default значение
	}

	claims := &FieldsClaims{
		email,
		strconv.Itoa(role),
		tokenType,
		jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   email,
			Issuer:    j.Issuer,
			Audience:  jwt.ClaimStrings{j.Audience},
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenString, err := token.SignedString(j.SecretKey)
	if err != nil {
		return "", fmt.Errorf("Ошибка при подписании токена: %v\n", err)
	}

	return tokenString, nil
}

func (j *JWT) GenerateRefreshToken(email string, role int) (string, error) {
	if email == "" {
		return "", fmt.Errorf("fields cannot be empty")
	}

	claims := &FieldsClaims{
		email,
		strconv.Itoa(role),
		"refresh",
		jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   email,
			Issuer:    j.Issuer,
			Audience:  jwt.ClaimStrings{j.Audience},
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenString, err := token.SignedString(j.SecretKey)
	if err != nil {
		return "", fmt.Errorf("Ошибка при подписании токена: %v\n", err)
	}
	return tokenString, nil
}

func (j *JWT) Verify(token string) (*FieldsClaims, error) {
	tokens, err := jwt.ParseWithClaims(token, &FieldsClaims{}, func(tok *jwt.Token) (interface{}, error) {
		if _, ok := tok.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", tok.Header["alg"])
		}
		return j.SecretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, ErrExpiredToken
		}

		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	if !tokens.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := tokens.Claims.(*FieldsClaims)
	if !ok {
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	return claims, nil
}

func (j *JWT) VerifyAccessToken(token string) (*FieldsClaims, error) {
	c, err := j.Verify(token)
	if err != nil {
		return nil, err
	}

	if c.Type != "access" {
		return nil, errors.New("invalid token type")
	}

	return c, nil
}

func (j *JWT) VerifyRefreshToken(token string) (*FieldsClaims, error) {
	c, err := j.Verify(token)
	if err != nil {
		return nil, err
	}

	if c.Type != "refresh" {
		return nil, errors.New("invalid token type")
	}

	return c, nil
}
