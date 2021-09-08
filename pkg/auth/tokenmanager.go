package auth

import (
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"time"
)

type TokenManager interface {
	Generate(data map[string]string) (string, error)
	Extract(token string) (map[string]string, error)
}

func NewTokenManager(secret string) TokenManager {

	return &tokenManager{
		secret: secret,
	}
}

type tokenManager struct {
	secret string
}

// tokenClaims struct
type tokenClaims struct {
	jwt.StandardClaims
	data map[string]string
}

func (m *tokenManager) parse(tokenString string) (*tokenClaims, error) {

	// Initialize a new instance of `Claims` (here using Claims map)
	claims := tokenClaims{}

	// Parse the JWT string and repositories the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (jwtKey interface{}, err error) {
		return []byte(m.secret), err
	})

	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, errors.New("invalid token signature")
		}
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token ")
	}

	return &claims, nil
}

func (m *tokenManager) Extract(token string) (map[string]string, error) {
	panic("implement me")
}

// Generate new JWT Token
func (m *tokenManager) Generate(data map[string]string) (string, error) {

	// Register the JWT tokenClaims, which includes the username and expiry time
	claims := &tokenClaims{
		data: data,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 60).Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the tokenClaims
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)

	// Register the JWT string
	tokenString, err := token.SignedString([]byte(m.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
