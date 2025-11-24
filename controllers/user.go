package controllers

import (
	"beego-jwt-vergo/models"
	"beego-jwt-vergo/services"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/astaxie/beego"
)

type UserController struct {
	beego.Controller
}

type AuthorizedResponse struct {
	Message string       `json:"message"`
	User    *models.User `json:"user"`
	Token   string       `json:"token"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

// ============================================================================
// REGISTER
// ============================================================================

// @router /register [post]
func (cont *UserController) RegisterUser() {

	fmt.Println("===== DEBUG REGISTER =====")

	// Lire le body correctement
	body := cont.Ctx.Input.CopyBody(1 << 20)

	fmt.Println("RAW BODY:", string(body))
	fmt.Println("Headers:", cont.Ctx.Request.Header)

	var iu models.InputUser
	if err := json.Unmarshal(body, &iu); err != nil {
		fmt.Println("JSON ERROR:", err)
		cont.Data["json"] = ErrorResponse{Message: "Invalid JSON format"}
		cont.ServeJSON()
		return
	}

	fmt.Println("PARSED USER:", iu)

	if iu.Email == "" || iu.Password == "" || iu.Name == "" {
		cont.Data["json"] = ErrorResponse{
			Message: "Missing fields: email, password, name",
		}
		cont.ServeJSON()
		return
	}

	id, err := models.CreateNew(iu.Email, iu.Password, iu.Name)
	if err != nil {
		cont.Data["json"] = ErrorResponse{Message: err.Error()}
		cont.ServeJSON()
		return
	}

	user, err := models.FindById(id)
	if err != nil {
		cont.Data["json"] = ErrorResponse{Message: err.Error()}
		cont.ServeJSON()
		return
	}

	token, err := services.MakeToken(id)
	if err != nil {
		cont.Data["json"] = ErrorResponse{Message: err.Error()}
		cont.ServeJSON()
		return
	}

	cont.Data["json"] = AuthorizedResponse{
		Message: "User created successfully",
		User:    user,
		Token:   token,
	}
	cont.ServeJSON()
}

// ============================================================================
// LOGIN
// ============================================================================

// @router /login [post]
func (cont *UserController) LoginUser() {

	body := cont.Ctx.Input.CopyBody(1 << 20)

	fmt.Println("===== DEBUG LOGIN =====")
	fmt.Println("RAW BODY:", string(body))

	var credentials models.BasicCredentials
	if err := json.Unmarshal(body, &credentials); err != nil {
		cont.Data["json"] = ErrorResponse{Message: "Invalid JSON format"}
		cont.ServeJSON()
		return
	}

	if credentials.Email == "" || credentials.Password == "" {
		cont.Data["json"] = ErrorResponse{Message: "Missing email or password"}
		cont.ServeJSON()
		return
	}

	user, err := models.Login(credentials.Email, credentials.Password)
	if err != nil {
		cont.Data["json"] = ErrorResponse{Message: err.Error()}
		cont.ServeJSON()
		return
	}

	token, err := services.MakeToken(user.Id)
	if err != nil {
		cont.Data["json"] = ErrorResponse{Message: err.Error()}
		cont.ServeJSON()
		return
	}

	cont.Data["json"] = AuthorizedResponse{
		Message: "User logged in successfully",
		User:    user,
		Token:   token,
	}
	cont.ServeJSON()
}

// @Title Index Users
// @Description Index all users when request is authorized
// @Param   Authorization   header   string   true   "Bearer token"
// @router /users [get]
func (cont *UserController) IndexAll() {

	fmt.Println("===== DEBUG USERS =====")

	// 1. Lire le header Authorization
	authHeader := cont.Ctx.Request.Header.Get("Authorization")
	fmt.Println("Authorization Header:", authHeader)

	if authHeader == "" {
		cont.Data["json"] = ErrorResponse{Message: "Missing Authorization header"}
		cont.ServeJSON()
		return
	}

	// Format attendu : "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		cont.Data["json"] = ErrorResponse{Message: "Invalid Authorization header format"}
		cont.ServeJSON()
		return
	}

	tokenString := parts[1]

	// 2. Extraire userId depuis le token RSA
	userId, err := services.GetUserIdFromToken(tokenString)
	if err != nil {
		fmt.Println("TOKEN ERROR:", err)
		cont.Data["json"] = ErrorResponse{Message: "Invalid or expired token"}
		cont.ServeJSON()
		return
	}

	fmt.Println("Token OK — userId =", userId)

	// 3. Récupérer les users dans MySQL
	users, err := models.IndexAll()
	if err != nil {
		cont.Data["json"] = ErrorResponse{Message: err.Error()}
		cont.ServeJSON()
		return
	}

	// 4. Retourner la liste des utilisateurs
	cont.Data["json"] = users
	cont.ServeJSON()
}
