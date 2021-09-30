package jwt

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/DCloudGaming/cloud-morph-host/pkg/env"
	"github.com/DCloudGaming/cloud-morph-host/pkg/errors"
	"github.com/DCloudGaming/cloud-morph-host/pkg/model"
	"github.com/DCloudGaming/cloud-morph-host/pkg/write"
	jwt "github.com/dgrijalva/jwt-go"
)

// jwt-cookie building and parsing
const cookieName = "pngr-jwt"

// tokens auto-refresh at the end of their lifetime,
// so long as the user hasn't been disabled in the interim
const tokenLifetime = time.Hour * 6

var hmacSecret []byte

func init() {
	hmacSecret = []byte(os.Getenv("TOKEN_SECRET"))
	if hmacSecret == nil {
		panic("No TOKEN_SECRET environment variable was found")
	}
}

type claims struct {
	User *model.User
	jwt.StandardClaims
}

// RequireAuth middleware makes sure the user exists based on their JWT
func RequireAuth(minStatus model.Status, e env.SharedEnv, w http.ResponseWriter, r *http.Request) (*model.User, bool) {
	u, err := HandleUserCookie(e, w, r)
	if err != nil {
		write.Error(err, w, r)
		return u, false
	}

	if u.Status < minStatus {
		write.Error(errors.RouteUnauthorized, w, r)
		return u, false
	}

	return u, true
}

// WriteUserCookie encodes a user's JWT and sets it as an httpOnly & Secure cookie
func WriteUserCookie(w http.ResponseWriter, u *model.User) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    encodeUser(u),
		Path:     "/",
		HttpOnly: true,
		//Secure:   true,
	})
}

//// HandleUserCookie attempts to refresh an expired token if the user is still valid
func HandleUserCookie(e env.SharedEnv, w http.ResponseWriter, r *http.Request) (*model.User, error) {
	u, err := userFromCookie(r)

	// attempt refresh of expired token:
	if err == errors.ExpiredToken && u.Status > 0 {
		user, fetchError := e.UserRepo().GetUser(u.WalletAddress)
		if fetchError != nil {
			return nil, err
		}
		if user.Status > 0 {
			WriteUserCookie(w, user)
			return user, nil
		}
	}

	return u, err
}

// userFromCookie builds a user object from a JWT, if it's valid
func userFromCookie(r *http.Request) (*model.User, error) {
	cookie, _ := r.Cookie(cookieName)
	var tokenString string
	if cookie != nil {
		tokenString = cookie.Value
	}

	if tokenString == "" {
		return &model.User{}, nil
	}

	return decodeUser(tokenString)
}

// encodeUser convert a user struct into a jwt
func encodeUser(u *model.User) (tokenString string) {
	claims := claims{
		u,
		jwt.StandardClaims{
			IssuedAt:  time.Now().Add(-time.Second).Unix(),
			ExpiresAt: time.Now().Add(tokenLifetime).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// unhandled err here
	tokenString, err := token.SignedString(hmacSecret)
	if err != nil {
		log.Println("Error signing token", err)
	}
	return
}

// decodeUser converts a jwt into a user struct (or returns a zero-value user)
func decodeUser(tokenString string) (*model.User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &claims{}, func(token *jwt.Token) (interface{}, error) {
		return hmacSecret, nil
	})

	if err != nil {
		// check for expired token
		if verr, ok := err.(*jwt.ValidationError); ok {
			if verr.Errors&jwt.ValidationErrorExpired != 0 {
				return getUserFromToken(token), errors.ExpiredToken
			}
		}
	}

	if err != nil || !token.Valid {
		return nil, errors.InvalidToken
	}

	return getUserFromToken(token), nil
}

func getUserFromToken(token *jwt.Token) *model.User {
	if claims, ok := token.Claims.(*claims); ok {
		return claims.User
	}

	return &model.User{}
}