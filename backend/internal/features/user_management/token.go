package usermanagement

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UID   string
	Email string
	Phone string
	jwt.RegisteredClaims
}

func GetClaims(uid, email, phone, role, issuer, audience string, duration time.Duration) (*Claims, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("tokenid generation failed: %v", err)
	}

	now := time.Now()

	return &Claims{
		UID:   uid,
		Email: email,
		Phone: phone,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID.String(),
			Issuer:    issuer,
			Audience:  jwt.ClaimStrings{audience},
			Subject:   uid,
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}, nil
}

type JWTManagerInterface interface {
	GenerateToken(uid, email, phone, role, issuer, audience string, duration time.Duration) (tokenStr string, err error)
	ValidateToken(tokenStr string) (*Claims, error)
}

type JwtManager struct {
	secret string
}

func NewJwtManager(secret string) JWTManagerInterface {
	return &JwtManager{
		secret: secret,
	}
}

func (j JwtManager) GenerateToken(uid, email, phone, role, issuer, audience string, duration time.Duration) (tokenStr string, err error) {
	claims, err := GetClaims(
		uid,
		email,
		phone,
		role,
		issuer,
		audience,
		duration,
	)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims, nil)
	tokenStr, err = token.SignedString(j.secret)

	return tokenStr, err
}
func (j JwtManager) ValidateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(t *jwt.Token) (any, error) {
			_, k := t.Method.(*jwt.SigningMethodRSA)
			if !k {
				return nil, ErrInvalidTokenSigningMethod
			}

			return []byte(j.secret), nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	claims, k := token.Claims.(*Claims)
	if !k {
		return nil, ErrInvalidClaimsFormat
	}

	if time.Now().Unix() > claims.ExpiresAt.Unix() {
		return nil, ErrTokenExpired
	}

	return claims, nil

}
