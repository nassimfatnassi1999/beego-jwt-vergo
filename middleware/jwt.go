package middleware

import (
    "beego-jwt-vergo/services"
    "net/http"
    "strings"

    "github.com/astaxie/beego"
    "github.com/astaxie/beego/context"
)

// RequireJWT is a Beego filter that ensures a valid JWT is present
func RequireJWT() beego.FilterFunc {
    return func(ctx *context.Context) {
        auth := ctx.Input.Header("Authorization")
        if auth == "" || !strings.HasPrefix(strings.ToLower(auth), "bearer ") {
            ctx.Output.SetStatus(http.StatusUnauthorized)
            _ = ctx.Output.JSON(map[string]string{"error": "missing or invalid Authorization header"}, false, false)
            return
        }

        token := strings.TrimSpace(auth[len("Bearer "):])

        if _, err := services.GetUserIdFromToken(token); err != nil {
            ctx.Output.SetStatus(http.StatusUnauthorized)
            _ = ctx.Output.JSON(map[string]string{"error": "invalid or expired token"}, false, false)
            return
        }
    }
}
