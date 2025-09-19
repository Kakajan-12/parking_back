package core

import (
	"backend/config"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// LoadRSAPrivateKey loads RSA private key from file
func LoadRSAPrivateKey(path string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(data)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// LoadRSAPublicKey loads RSA public key from file
func LoadRSAPublicKey(path string) (*rsa.PublicKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA public key")
	}

	return rsaPub, nil
}

// JWTPayload creates claims
func JWTPayload(
	data map[string]interface{},
	iat *time.Time,
	expireDelta *time.Duration,
) jwt.MapClaims {
	issuedAt := iat
	if issuedAt == nil {
		t := time.Now()
		issuedAt = &t
	}

	var expire time.Time
	if expireDelta != nil {
		expire = (*issuedAt).Add(*expireDelta)
	} else {
		expire = (*issuedAt).Add(24 * time.Hour)
	}

	claims := jwt.MapClaims{}
	for k, v := range data {
		claims[k] = v
	}

	claims["iat"] = (*issuedAt).Unix()
	claims["exp"] = expire.Unix()
	claims["iss"] = "backend"
	claims["aud"] = "client"

	return claims
}

// JWTEncode signs claims using RS256
func JWTEncode(claims jwt.Claims) (string, error) {
	keyPath := config.AppConfig.JWTPrivateKeyPath
	privateKey, err := LoadRSAPrivateKey(keyPath)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}

// JWTDecode parses and verifies RS256 token
func JWTDecode(tokenStr string) (jwt.MapClaims, error) {
	keyPath := config.AppConfig.JWTPublicKeyPath
	publicKey, err := LoadRSAPublicKey(keyPath)
	if err != nil {
		return nil, err
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != "RS256" {
			return nil, errors.New("unexpected signing method")
		}
		return publicKey, nil
	})

	if err != nil {
		var ve *jwt.ValidationError
		if errors.As(err, &ve) && ve.Errors&jwt.ValidationErrorExpired != 0 {
			return nil, errors.New("token is expired")
		}
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	if !claims.VerifyAudience(config.AppConfig.JwtAudience, true) {
		return nil, errors.New("invalid audience")
	}

	if !claims.VerifyIssuer(config.AppConfig.JwtIssuer, true) {
		return nil, errors.New("invalid issuer")
	}

	return claims, nil
}
