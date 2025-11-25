package middleware

import (
    "beego-jwt-vergo/services"
    "net/http"
    "strings"

    "github.com/astaxie/beego/context"
)

// AdminOnly ensures the JWT belongs to an admin
func AdminOnly(ctx *context.Context) {
    auth := ctx.Input.Header("Authorization")
    if auth == "" || !strings.HasPrefix(strings.ToLower(auth), "bearer ") {
        ctx.Output.SetStatus(http.StatusForbidden)
        _ = ctx.Output.JSON(map[string]string{"error": "admin token required"}, false, false)
        return
    }

    token := strings.TrimSpace(auth[len("Bearer "):])

    role, err := services.GetRoleFromToken(token)
    if err != nil {
        ctx.Output.SetStatus(http.StatusForbidden)
        _ = ctx.Output.JSON(map[string]string{"error": "invalid token"}, false, false)
        return
    }

    if role != "admin" {
        ctx.Output.SetStatus(http.StatusForbidden)
        _ = ctx.Output.JSON(map[string]string{"error": "admin access required"}, false, false)
        return
    }
}
