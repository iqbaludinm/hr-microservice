package domain

import (
	"fmt"
	"strconv"
	"time"

	"github.com/iqbaludinm/hr-microservice/profile-service/helper"
	"github.com/iqbaludinm/hr-microservice/profile-service/model/web"
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

func (s *User) ToUserResponse() web.UserResponse {
	return web.UserResponse{
		ID:    s.ID,
		Name:  s.Name,
		Email: s.Email,
		Phone: s.Phone,
	}
}

// Helper function for converting the JobQueryFilter from web to domain
// This function is used for calling a function in the 'repository' layer
func ToDomainUserQueryFilter(q web.UserQueryFilter) UserQueryFilter {
	var showDeleted bool
	if q.ShowDeleted == true {
		showDeleted = true
	}

	// pagination validation. If the value is not a number, then we will use the default value.
	var page, limit string
	if q.Page != "" {
		pageInt, _ := strconv.Atoi(q.Page)
		if q.Limit != "" {
			limitInt, _ := strconv.Atoi(q.Limit)
			page = fmt.Sprintf("%d", (pageInt-1)*limitInt)
			limit = fmt.Sprintf("%d", limitInt)
		} else {
			page = fmt.Sprintf("%d", (pageInt-1)*helper.DefaultLimit)
			limit = fmt.Sprintf("%d", helper.DefaultLimit)
		}
	} else {
		if q.Limit != "" {
			limitInt, _ := strconv.Atoi(q.Limit)
			page = "0"
			limit = fmt.Sprintf("%d", limitInt)
		}
	}

	return UserQueryFilter{
		Name:        q.Name,
		Email:       q.Email,
		Phone:       q.Phone,
		ShowDeleted: showDeleted,
		Pagination: Pagination{
			Page:  page,
			Limit: limit,
		},
	}
}
