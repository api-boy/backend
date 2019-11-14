package authutils

import (
	"fmt"
	"time"

	"apiboy/backend/src/errors"

	jwtgo "github.com/dgrijalva/jwt-go"
)

// jwtClaims contains the JWT custom claims for the app
type jwtClaims struct {
	jwtgo.StandardClaims
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`
	UserRole  string `json:"user_role"`
}

// NewJWT returns a new JWT token
func NewJWT(jwtIssuer, jwtSignKey string, data *AuthData) (string, error) {
	// create token
	now := time.Now().UTC()
	expiration := now.Add(30 * 24 * time.Hour).UTC() // 30 days

	claims := jwtClaims{
		jwtgo.StandardClaims{
			Id:        data.JwtID,
			Subject:   fmt.Sprintf("%v", data.UserID),
			Issuer:    jwtIssuer,
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
			ExpiresAt: expiration.Unix(),
		},
		data.UserName,
		data.UserEmail,
		data.UserRole,
	}

	token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, &claims)

	// sign token
	jwt, err := token.SignedString([]byte(jwtSignKey))
	if err != nil {
		return "", err
	}

	return jwt, nil
}

// ParseJWT parses a JWT token
func ParseJWT(jwtSignKey, jwt string) (*AuthData, error) {
	// parse token
	token, err := jwtgo.ParseWithClaims(jwt, &jwtClaims{}, func(token *jwtgo.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtgo.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSignKey), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.Unauthenticated{Msg: "Could not parse jwt token", Err: err}
	}

	// get claims
	claims, ok := token.Claims.(*jwtClaims)
	if !ok {
		return nil, errors.Unauthenticated{Msg: "Could not get token claims"}
	}

	return &AuthData{
		JwtID:     claims.Id,
		UserID:    claims.Subject,
		UserName:  claims.UserName,
		UserEmail: claims.UserEmail,
		UserRole:  claims.UserRole,
	}, nil
}
