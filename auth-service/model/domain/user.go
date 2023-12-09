package domain

import (
	"time"

	"github.com/iqbaludinm/hr-microservice/auth-service/model/web"
	"golang.org/x/crypto/bcrypt"
)

// user main struct
type User struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	Phone     string     `json:"phone"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (user *User) SetPassword(password string) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	user.Password = string(hashedPassword)
}

func (user *User) ComparePassword(correctPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(correctPassword), []byte(password))
}

func ToRegisterResponse(user User) web.RegisterResponse {
	return web.RegisterResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
	}
}

func ToLoginResponse(user User) web.LoginResponse {

	return web.LoginResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
	}
}

// type UserWithName struct {
// 	Id              *int       `json:"id"`
// 	Name            string     `json:"name"`
// 	Email           string     `json:"email"`
// 	Password        string     `json:"password"`
// 	Phone           string     `json:"phone"`
// 	CreatedAt       time.Time  `json:"created_at"`
// 	UpdatedAt       time.Time  `json:"updated_at"`
// 	DeletedAt       *time.Time `json:"deleted_at"`
// }
