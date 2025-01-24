package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/LikheKeto/Suraksheet/config"
	"github.com/LikheKeto/Suraksheet/types"
	"github.com/LikheKeto/Suraksheet/utils"
	"github.com/golang-jwt/jwt/v5"
)

var user = new(types.User)

func CreateJWT(secret []byte, userID int) (string, error) {
	expiration := time.Hour * 24 * 30
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

func ValidateJWT(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(config.Envs.JWTSecret), nil
	})
}

func WithJWTAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the token from Authentication headers
		split := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		tokenString := split[len(split)-1]

		token, err := ValidateJWT(tokenString)
		if err != nil {
			log.Printf("failed to validate token: %v", err)
			permissionDenied(w)
			return
		}

		// validate the JWT
		if !token.Valid {
			log.Println("invalid token")
			permissionDenied(w)
			return
		}

		// fetch userID from DB
		claims := token.Claims.(jwt.MapClaims)
		str, ok := claims["userID"].(string)
		if !ok {
			permissionDenied(w)
			return
		}

		userID, err := strconv.Atoi(str)
		if err != nil {
			permissionDenied(w)
			return
		}

		u, err := store.GetUserByID(userID)
		if err != nil {
			permissionDenied(w)
			return
		}

		// set context "userID" value
		ctx := r.Context()
		ctx = context.WithValue(ctx, user, u)
		r = r.WithContext(ctx)

		handlerFunc(w, r)
	}
}

func ExtractUserFromContext(r *http.Request) (*types.User, error) {
	usr := r.Context().Value(user)
	if usr == nil {
		return nil, fmt.Errorf("permission denied")
	}
	if user, ok := usr.(*types.User); !ok {
		return nil, fmt.Errorf("permission denied")
	} else {
		return user, nil
	}
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}
