package services

import (
    "errors"
    "time"

    "github.com/astaxie/beego"
    jwt "github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

func init() {
    secret := beego.AppConfig.String("jwt_secret")
    if secret == "" {
        // Fallback dev secret; DO NOT use in production
        secret = "dev-secret-change-me"
    }
    jwtSecret = []byte(secret)
}

type Claims struct {
    UserID int64  `json:"uid"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

// MakeToken creates a signed JWT for a user
func MakeToken(userID int64, role string) (string, error) {
    if role == "" {
        role = "user"
    }

    claims := &Claims{
        UserID: userID,
        Role:   role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

func parseToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return jwtSecret, nil
    })
    if err != nil {
        return nil, err
    }

    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        return nil, errors.New("invalid token")
    }

    return claims, nil
}

// GetUserIdFromToken extracts user id from a token
func GetUserIdFromToken(tokenString string) (int64, error) {
    claims, err := parseToken(tokenString)
    if err != nil {
        return 0, err
    }
    return claims.UserID, nil
}

// GetRoleFromToken extracts role from a token
func GetRoleFromToken(tokenString string) (string, error) {
    claims, err := parseToken(tokenString)
    if err != nil {
        return "", err
    }
    return claims.Role, nil
}
