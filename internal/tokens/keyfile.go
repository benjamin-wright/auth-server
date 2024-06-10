package tokens

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var Issuer = "PongleHub"

type Keyfile []byte

type Claims struct {
	Subject   string
	Audiences []string
}

func (k *Keyfile) New(id string, audiences []string, expiration time.Duration) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Subject:   id,
			Issuer:    Issuer,
			Audience:  audiences,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
		},
	)

	key := []byte(*k)

	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %+v", err)
	}

	return tokenString, nil
}

func (k *Keyfile) Parse(token string) (Claims, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(*k), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return Claims{}, fmt.Errorf("couldn't parse non-token object: %s", token)
		} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return Claims{}, fmt.Errorf("invalid signature: %s", t.Signature)
		} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
			return Claims{}, errors.New("token expired or not ready yet")
		} else {
			return Claims{}, fmt.Errorf("error parsing token: %+v", err)
		}
	}

	issuer, err := t.Claims.GetIssuer()
	if err != nil {
		return Claims{}, fmt.Errorf("error getting issuer: %+v", err)
	}

	if issuer != Issuer {
		return Claims{}, fmt.Errorf("invalid issuer, expected PongleHub, got %s", issuer)
	}

	audience, err := t.Claims.GetAudience()
	if err != nil {
		return Claims{}, fmt.Errorf("error getting audience: %+v", err)
	}

	subject, err := t.Claims.GetSubject()
	if err != nil {
		return Claims{}, fmt.Errorf("error getting subject: %+v", err)
	}

	return Claims{
		Subject:   subject,
		Audiences: audience,
	}, nil
}
