package api

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ignoxx/podara/poc3/types"
	log "github.com/sirupsen/logrus"
)

type UserClaims struct {
	Email string `json:"email"`
	Id    string `json:"id"`
	jwt.MapClaims
}

func withAuth(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		receivedToken := r.Header.Get("Authorization")

		if receivedToken == "" {
			// get token from cookie
			cookie, err := r.Cookie("token")

			if err != nil {
				log.WithFields(log.Fields{
					"method": r.Method,
					"path":   r.URL.Path,
					"addr":   r.RemoteAddr,
				}).Errorf("%s", err.Error())

				WriteJSON(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
				return
			}

			receivedToken = cookie.Value
		}

		token, err := validateJWT(receivedToken)

		if err != nil {
			log.WithFields(log.Fields{
				"method": r.Method,
				"path":   r.URL.Path,
				"addr":   r.RemoteAddr,
			}).Errorf("%s", err.Error())

			WriteJSON(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			return
		}

		if !token.Valid {
			log.WithFields(log.Fields{
				"method": r.Method,
				"path":   r.URL.Path,
				"addr":   r.RemoteAddr,
			}).Errorf("%s", err.Error())

			WriteJSON(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			return
		}

		f(w, r)
	}
}

func getJwtClaims(r *http.Request) (*UserClaims, error) {
	receivedToken := r.Header.Get("Authorization")
	if receivedToken == "" {
		// get token from Cookie
		cookie, err := r.Cookie("token")
		if err != nil {
			return nil, err
		}
		receivedToken = cookie.Value
	}

	token, err := validateJWT(receivedToken)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*UserClaims)

	if !ok {
		return nil, fmt.Errorf("invalid claims type, %v, %v, %v", claims, ok, token.Claims)
	}

	return claims, nil
}

func createJwtToken(user *types.User) (string, error) {
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")

	claims := UserClaims{
		user.Email,
		user.Id,
		jwt.MapClaims{
			"exp": jwt.NewNumericDate(time.Now().Add(time.Hour)),
			"iss": "podara",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(jwtSecretKey))
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")

	return jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (any, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(jwtSecretKey), nil
	})
}

func setAuthCookie(w http.ResponseWriter, token string) {
	cookie := http.Cookie{
		Value:    token,
		Name:     "token",
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &cookie)
}
