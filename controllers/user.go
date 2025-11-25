package controllers

import (
    "beego-jwt-vergo/models"
    "beego-jwt-vergo/services"
    "encoding/json"
    "fmt"
    "net/http"
    "strings"

    "github.com/astaxie/beego"
)

type UserController struct {
    beego.Controller
}

type ErrorResponse struct {
    Error string `json:"error"`
}

type AuthorizedResponse struct {
    Message string       `json:"message"`
    User    *models.User `json:"user"`
    Token   string       `json:"token"`
}

// RegisterUser handles POST /register
func (c *UserController) RegisterUser() {
    body := c.Ctx.Input.CopyBody(1 << 20)

    var input models.InputUser
    if err := json.Unmarshal(body, &input); err != nil {
        c.CustomAbort(http.StatusBadRequest, "invalid JSON payload")
        return
    }

    if strings.TrimSpace(input.Email) == "" ||
        strings.TrimSpace(input.Password) == "" ||
        strings.TrimSpace(input.Name) == "" {
        c.Ctx.Output.SetStatus(http.StatusBadRequest)
        _ = c.Ctx.Output.JSON(ErrorResponse{"email, password and name are required"}, false, false)
        return
    }

    // Default role = user
    role := input.Role
    if role == "" {
        role = "user"
    }

    id, err := models.CreateNew(input.Email, input.Password, input.Name, role)
    if err != nil {
        c.Ctx.Output.SetStatus(http.StatusBadRequest)
        _ = c.Ctx.Output.JSON(ErrorResponse{err.Error()}, false, false)
        return
    }

    user, err := models.FindById(id)
    if err != nil {
        c.Ctx.Output.SetStatus(http.StatusInternalServerError)
        _ = c.Ctx.Output.JSON(ErrorResponse{err.Error()}, false, false)
        return
    }

    token, err := services.MakeToken(user.Id, user.Role)
    if err != nil {
        c.Ctx.Output.SetStatus(http.StatusInternalServerError)
        _ = c.Ctx.Output.JSON(ErrorResponse{err.Error()}, false, false)
        return
    }

    resp := AuthorizedResponse{
        Message: "user created successfully",
        User:    user,
        Token:   token,
    }
    _ = c.Ctx.Output.JSON(resp, false, false)
}

// LoginUser handles POST /login
func (c *UserController) LoginUser() {
    body := c.Ctx.Input.CopyBody(1 << 20)

    var creds models.BasicCredentials
    if err := json.Unmarshal(body, &creds); err != nil {
        c.Ctx.Output.SetStatus(http.StatusBadRequest)
        _ = c.Ctx.Output.JSON(ErrorResponse{"invalid JSON payload"}, false, false)
        return
    }

    if strings.TrimSpace(creds.Email) == "" || strings.TrimSpace(creds.Password) == "" {
        c.Ctx.Output.SetStatus(http.StatusBadRequest)
        _ = c.Ctx.Output.JSON(ErrorResponse{"email and password are required"}, false, false)
        return
    }

    user, err := models.Login(creds.Email, creds.Password)
    if err != nil {
        c.Ctx.Output.SetStatus(http.StatusUnauthorized)
        _ = c.Ctx.Output.JSON(ErrorResponse{err.Error()}, false, false)
        return
    }

    token, err := services.MakeToken(user.Id, user.Role)
    if err != nil {
        c.Ctx.Output.SetStatus(http.StatusInternalServerError)
        _ = c.Ctx.Output.JSON(ErrorResponse{err.Error()}, false, false)
        return
    }

    resp := AuthorizedResponse{
        Message: "user logged in successfully",
        User:    user,
        Token:   token,
    }
    _ = c.Ctx.Output.JSON(resp, false, false)
}

// IndexAll handles GET /users and /v1/users (admin only)
func (c *UserController) IndexAll() {
    authHeader := c.Ctx.Request.Header.Get("Authorization")
    if authHeader == "" {
        c.Ctx.Output.SetStatus(http.StatusUnauthorized)
        _ = c.Ctx.Output.JSON(ErrorResponse{"missing Authorization header"}, false, false)
        return
    }

    parts := strings.SplitN(authHeader, " ", 2)
    if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
        c.Ctx.Output.SetStatus(http.StatusUnauthorized)
        _ = c.Ctx.Output.JSON(ErrorResponse{"invalid Authorization header format"}, false, false)
        return
    }

    tokenString := strings.TrimSpace(parts[1])

    userID, err := services.GetUserIdFromToken(tokenString)
    if err != nil {
        c.Ctx.Output.SetStatus(http.StatusUnauthorized)
        _ = c.Ctx.Output.JSON(ErrorResponse{"invalid or expired token"}, false, false)
        return
    }

    role, err := services.GetRoleFromToken(tokenString)
    if err != nil {
        c.Ctx.Output.SetStatus(http.StatusForbidden)
        _ = c.Ctx.Output.JSON(ErrorResponse{"cannot read role from token"}, false, false)
        return
    }

    fmt.Println("Token OK — userId =", userID, "role =", role)

    if role != "admin" {
        c.Ctx.Output.SetStatus(http.StatusForbidden)
        _ = c.Ctx.Output.JSON(ErrorResponse{"admin access required"}, false, false)
        return
    }

    users, err := models.IndexAll()
    if err != nil {
        c.Ctx.Output.SetStatus(http.StatusInternalServerError)
        _ = c.Ctx.Output.JSON(ErrorResponse{err.Error()}, false, false)
        return
    }

    _ = c.Ctx.Output.JSON(users, false, false)
}
