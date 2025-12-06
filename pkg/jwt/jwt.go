package jwt

import (
	"errors"
	"fmt"
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
	Role uint   `json:"role"`
	Type string `json:"type"`
	jwt.RegisteredClaims
}

func New(secret []byte) *JWT {
	return &JWT{
		SecretKey: secret,
		Issuer:    "crm-system",
		Audience:  "zerat",
	}
}

type ReturnDataToken struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (j *JWT) GenerateAccessToken(email string, role uint, tokenType string) (*ReturnDataToken, error) {

	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}

	if tokenType == "" {
		tokenType = "access" // default значение
	}

	now := time.Now()
	expertise := now.Add(45 * time.Minute) // 45 min - balanced security for personal data

	claims := &FieldsClaims{
		role,
		tokenType,
		jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			ExpiresAt: jwt.NewNumericDate(expertise),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   email,
			Issuer:    j.Issuer,
			Audience:  jwt.ClaimStrings{j.Audience},
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenString, err := token.SignedString(j.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("Ошибка при подписании токена: %v\n", err)
	}

	return &ReturnDataToken{tokenString, expertise}, nil
}

func (j *JWT) GenerateRefreshToken(email string, role uint) (*ReturnDataToken, error) {
	if email == "" {
		return nil, fmt.Errorf("fields cannot be empty")
	}

	now := time.Now()
	expertise := now.Add(30 * 24 * time.Hour)

	claims := &FieldsClaims{
		role,
		"refresh",
		jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			ExpiresAt: jwt.NewNumericDate(expertise),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   email,
			Issuer:    j.Issuer,
			Audience:  jwt.ClaimStrings{j.Audience},
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenString, err := token.SignedString(j.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("Ошибка при подписании токена: %v\n", err)
	}
	return &ReturnDataToken{tokenString, expertise}, nil
}

func (j *JWT) Verify(token string) (*FieldsClaims, error) {
	t, err := jwt.ParseWithClaims(token, &FieldsClaims{}, func(tok *jwt.Token) (interface{}, error) {
		if _, ok := tok.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", tok.Header["alg"])
		}
		return j.SecretKey, nil
	})

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired), errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, ErrExpiredToken
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			return nil, ErrInvalidToken // или ErrInvalidSignature
		default:
			return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
		}
	}

	if !t.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := t.Claims.(*FieldsClaims)
	if !ok {
		return nil, ErrInvalidToken
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
