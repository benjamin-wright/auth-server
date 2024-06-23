package tokens

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type TokenVerifier struct {
	public *rsa.PublicKey
}

func NewVerifier(certPath string) (*TokenVerifier, error) {
	cert, err := os.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read cert file: %+v", err)
	}

	block, _ := pem.Decode(cert)
	if block == nil {
		return nil, errors.New("failed to decode public key")
	}

	parsed, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %+v", err)
	}

	key, ok := parsed.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("wrong kind of public key")
	}

	return &TokenVerifier{
		public: key,
	}, nil
}

type Claims struct {
	Subject   string
	Audiences []string
}

func (t *TokenVerifier) Parse(token string) (Claims, error) {
	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return t.public, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return Claims{}, fmt.Errorf("couldn't parse non-token object: %s", token)
		} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return Claims{}, fmt.Errorf("invalid signature: %s", parsed.Signature)
		} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
			return Claims{}, errors.New("token expired or not ready yet")
		} else {
			return Claims{}, fmt.Errorf("error parsing token: %+v", err)
		}
	}

	issuer, err := parsed.Claims.GetIssuer()
	if err != nil {
		return Claims{}, fmt.Errorf("error getting issuer: %+v", err)
	}

	if issuer != Issuer {
		return Claims{}, fmt.Errorf("invalid issuer, expected PongleHub, got %s", issuer)
	}

	audience, err := parsed.Claims.GetAudience()
	if err != nil {
		return Claims{}, fmt.Errorf("error getting audience: %+v", err)
	}

	subject, err := parsed.Claims.GetSubject()
	if err != nil {
		return Claims{}, fmt.Errorf("error getting subject: %+v", err)
	}

	return Claims{
		Subject:   subject,
		Audiences: audience,
	}, nil
}
