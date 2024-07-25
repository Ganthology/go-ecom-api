package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ganthology/go-ecom-api/config"
	"github.com/ganthology/go-ecom-api/types"
	"github.com/ganthology/go-ecom-api/utils"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserKey contextKey = "userID"

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

func WithJWTAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get token from request
		tokenString := getTokenFromRequest(r)
		// validate jwt
		token, err := validateToken(tokenString)
		if err != nil {
			log.Printf("error validating token: %v", err)
			permissionDenied(w)
			return
		}

		if !token.Valid {
			log.Printf("token is invalid")
			permissionDenied(w)
			return
		}
		// fetch userId from db
		claims := token.Claims.(jwt.MapClaims)
		str := claims["userID"].(string)

		userID, _ := strconv.Atoi(str)

		u, err := store.GetUserById(userID)
		if err != nil {
			log.Printf("error fetching user: %v", err)
			permissionDenied(w)
			return
		}
		// set context userID to user id
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, u.ID)
		r = r.WithContext(ctx)

		handlerFunc(w, r)
	}
}

func getTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	if tokenAuth == "" {
		return tokenAuth
	}
	return ""
}

func validateToken(tokenString string) (*jwt.Token, error) {
	secret := []byte(config.Envs.JWTSecret)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}

func GetUserIDFromContext(ctx context.Context) int {
	userID, ok := ctx.Value(UserKey).(int)
	if !ok {
		return -1
	}
	return userID
}
