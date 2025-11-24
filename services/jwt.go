package services

import (
	"crypto/rsa"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/SermoDigital/jose/crypto"
	jwt "github.com/golang-jwt/jwt/v5"
)

var (
	iss = "api"
	sub = "uid-"
	aud = "client"
	exp = 24 * 30 * time.Hour
	nbf = 30 * time.Second
)

type Claims struct {
	jwt.RegisteredClaims
}

func GetKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {

	privBytes, err := os.ReadFile("./keys/private.txt")
	if err != nil {
		return nil, nil, err
	}

	privKey, err := crypto.ParseRSAPrivateKeyFromPEM(privBytes)
	if err != nil {
		return nil, nil, err
	}

	pubBytes, err := os.ReadFile("./keys/public.txt")
	if err != nil {
		return nil, nil, err
	}

	pubKey, err := crypto.ParseRSAPublicKeyFromPEM(pubBytes)
	if err != nil {
		return nil, nil, err
	}

	return privKey, pubKey, nil
}

func MakeToken(uid int64) (string, error) {
	privKey, _, err := GetKeyPair()
	if err != nil {
		return "", err
	}

	now := time.Now()
	uidStr := strconv.FormatInt(uid, 10)

	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    iss,
			Subject:   sub + uidStr,
			Audience:  jwt.ClaimStrings{aud},
			ExpiresAt: jwt.NewNumericDate(now.Add(exp)),
			NotBefore: jwt.NewNumericDate(now.Add(nbf)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	return token.SignedString(privKey)
}

func ValidateToken(tokenString string, uid int64) (bool, error) {

	_, pubKey, err := GetKeyPair()
	if err != nil {
		return false, err
	}

	claims := &Claims{}

	// Vérifie signature + claims
	_, err = jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return pubKey, nil
	})

	if err != nil {
		return false, err
	}

	uidStr := strconv.FormatInt(uid, 10)

	if claims.Subject != sub+uidStr {
		return false, nil
	}

	return true, nil
}

// GetUserIdFromToken extrait l'ID utilisateur (uid) depuis le token JWT.
func GetUserIdFromToken(tokenString string) (int64, error) {
	// On a juste besoin de la clé publique pour vérifier le token
	_, pubKey, err := GetKeyPair()
	if err != nil {
		return 0, err
	}

	claims := &Claims{}

	// Vérifie la signature + remplit les claims
	_, err = jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return pubKey, nil
	})
	if err != nil {
		return 0, err
	}

	// Subject attendu : "uid-<id>"
	if !strings.HasPrefix(claims.Subject, sub) {
		return 0, errors.New("invalid subject in token")
	}

	idStr := strings.TrimPrefix(claims.Subject, sub)
	uid, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, errors.New("invalid user id in token")
	}

	return uid, nil
}
