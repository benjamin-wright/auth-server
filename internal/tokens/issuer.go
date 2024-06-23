package tokens

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var Issuer = "PongleHub"

type TokenIssuer struct {
	private *rsa.PrivateKey
}

func NewIssuer(keyPath string) (*TokenIssuer, error) {
	key, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %+v", err)
	}

	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.New("failed to decode private key")
	}

	parsedKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %+v", err)
	}

	return &TokenIssuer{
		private: parsedKey,
	}, nil
}

func (t *TokenIssuer) New(id string, audiences []string, expiration time.Duration) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodRS512,
		jwt.RegisteredClaims{
			Subject:   id,
			Issuer:    Issuer,
			Audience:  audiences,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
		},
	)

	tokenString, err := token.SignedString(t.private)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %+v", err)
	}

	return tokenString, nil
}
