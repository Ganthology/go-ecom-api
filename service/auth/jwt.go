package auth

import (
	"strconv"
	"time"

	"github.com/ganthology/go-ecom-api/config"
	"github.com/golang-jwt/jwt/v5"
)

func CreateJWT(secret []byte, userID int) (string, error) {
	expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)

	// SigningMethodHS256 is a signing method using HMAC and SHA-256,
	// Symmetric key algorithm,
	// Single secret key for both signing and verification
	// --
	// SigningMethodES256 is a signing method using ECDSA and P-256 + SHA-256,
	// Asymmetric key algorithm,
	// Public/private key pair for signing and verification
	// --
	// privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	// token := jwt.New(jwt.SigningMethodES256)
	// tokenString, err := token.SignedString(privateKey)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(userID),
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
