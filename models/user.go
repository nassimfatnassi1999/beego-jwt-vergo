package models

import (
    "errors"
    "strings"
    "time"

    "github.com/astaxie/beego/orm"
    "golang.org/x/crypto/bcrypt"
)

// User represents a user in database and in JSON
type User struct {
    Id        int64     `json:"id" orm:"pk;auto"`
    Email     string    `json:"email" orm:"unique;index;size(191)"`
    Password  string    `json:"-" orm:"size(255)"`
    Name      string    `json:"name" orm:"size(100)"`
    Role      string    `json:"role" orm:"size(20);default(user)"`
    CreatedAt time.Time `json:"created_on" orm:"auto_now_add;type(datetime)"`
    UpdatedAt time.Time `json:"updated_on" orm:"auto_now;type(datetime)"`
}

// InputUser is used when registering
type InputUser struct {
    Email    string `json:"email"`
    Password string `json:"password"`
    Name     string `json:"name"`
    Role     string `json:"role"`
}

// BasicCredentials is used for login
type BasicCredentials struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

func init() {
    orm.RegisterModel(new(User))
}

func (u *User) TableName() string {
    return "users"
}

// IndexAll returns all users
func IndexAll() ([]*User, error) {
    o := orm.NewOrm()
    var users []*User
    _, err := o.QueryTable(new(User)).All(&users)
    return users, err
}

// FindById returns a user by id
func FindById(id int64) (*User, error) {
    o := orm.NewOrm()
    user := &User{Id: id}
    if err := o.Read(user); err != nil {
        if err == orm.ErrNoRows {
            return nil, errors.New("user not found")
        }
        return nil, err
    }
    return user, nil
}

// FindByEmail returns a user by e-mail
func FindByEmail(email string) (*User, error) {
    o := orm.NewOrm()
    user := &User{}
    err := o.QueryTable(new(User)).Filter("Email", email).One(user)
    if err != nil {
        if err == orm.ErrNoRows {
            return nil, errors.New("user not found")
        }
        return nil, err
    }
    return user, nil
}

// CreateNew creates a new user with hashed password
func CreateNew(email, password, name, role string) (int64, error) {
    email = strings.TrimSpace(email)
    password = strings.TrimSpace(password)
    name = strings.TrimSpace(name)
    role = strings.TrimSpace(role)

    if email == "" || password == "" || name == "" {
        return 0, errors.New("email, password and name are required")
    }

    if role == "" {
        role = "user"
    }

    // Check if email already exists
    if _, err := FindByEmail(email); err == nil {
        return 0, errors.New("email already exists")
    }

    // Hash password
    passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return 0, err
    }

    u := &User{
        Email:    email,
        Password: string(passHash),
        Name:     name,
        Role:     role,
    }

    o := orm.NewOrm()
    id, err := o.Insert(u)
    if err != nil {
        return 0, err
    }
    return id, nil
}

// Login validates user credentials and returns the user
func Login(email, password string) (*User, error) {
    email = strings.TrimSpace(email)
    password = strings.TrimSpace(password)

    if email == "" || password == "" {
        return nil, errors.New("email and password are required")
    }

    u, err := FindByEmail(email)
    if err != nil {
        return nil, errors.New("invalid email or password")
    }

    if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
        return nil, errors.New("invalid email or password")
    }

    return u, nil
}
