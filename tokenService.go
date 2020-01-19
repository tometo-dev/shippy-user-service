package main

import (
	"github.com/dgrijalva/jwt-go"
	pb "github.com/tsuki42/shippy-user-service/proto/auth"
	"time"
)

var (
	// Define a secure key string
	// used as a salt when hashing our tokens.
	key = []byte("mySuperSecretKeyLol")
)

// CustomClaims is our custom metadata, which will be hashed and send
// as the second segment in our JWT
type CustomClaims struct {
	User *pb.User
	jwt.StandardClaims
}

type Authable interface {
	Decode(token string) (*CustomClaims, error)
	Encode(user *pb.User) (string, error)
}

type TokenService struct {
	repo Repository
}

// Decode a token string into a token object
func (srv *TokenService) Decode(tokenString string) (*CustomClaims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	// Validate the token and return the custom claims
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

// Encode a claim into a JWT
func (srv *TokenService) Encode(user *pb.User) (string, error) {
	expireTime := time.Now().Add(time.Hour * 72).Unix()

	// Create the claims
	claims := CustomClaims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime,
			Issuer:    "shippy.user",
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	// Sign token and return
	return token.SignedString(key)
}
