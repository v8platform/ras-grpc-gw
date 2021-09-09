package auth

import (
	"bytes"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"

	"time"
)

type TokenManager interface {
	Generate(issuer, subject string, ttl time.Duration) (string, error)
	Validate(payload string, subject string) (string, error)
}

func NewTokenManager(secret string) (TokenManager, error) {
	digest := sha256.Sum256([]byte(secret))

	privKey, err := rsa.GenerateKey(bytes.NewBuffer(digest[:]), 64)
	if err != nil {
		return nil, err
	}

	return &tokenManager{
		key: privKey,
	}, nil
}

type tokenManager struct {
	key *rsa.PrivateKey
}

func (m *tokenManager) Validate(payload string, subject string) (string, error) {
	token, err := jwt.Parse(
		[]byte(payload),
		jwt.WithValidate(true),
		jwt.WithVerify(jwa.RS256, m.key.PublicKey),
		jwt.WithSubject(subject),
		// jwt.WithIssuer(issuer),
	)
	if err != nil {
		return "", err
	}
	return token.Issuer(), nil
}

// Generate new JWT Token
func (m *tokenManager) Generate(issuer, subject string, ttl time.Duration) (string, error) {

	t := jwt.New()
	t.Set(jwt.SubjectKey, subject)
	t.Set(jwt.IssuerKey, issuer)
	t.Set(jwt.ExpirationKey, time.Now().Add(ttl))
	t.Set(jwt.IssuedAtKey, time.Now())

	payload, err := jwt.Sign(t, jwa.RS256, m.key)
	if err != nil {
		fmt.Printf("failed to generate signed payload: %s\n", err)
		return "", err
	}

	return string(payload), nil
}
