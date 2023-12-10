package web

import "time"

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min:6"`
	Phone    string `json:"phone" validate:"required"`
}

type UpdateProfileRequest struct {
	Name      string `json:"name"`
	Email     string `json:"email" validate:"email"`
	Phone     string `json:"phone" `
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateUserPasswordRequest struct {
	Password string `json:"password" validate:"min:6"`
}

type UserQueryFilter struct {
	Name     string `query:"name"`
	Email    string `query:"email"`
	Phone    string `query:"phone"`
	Password string `query:"password"`

	// ShowDeleted is used for showing the soft-deleted user or not.
	ShowDeleted bool `query:"show_deleted"`

	// Pagination is used for fetching the data by page. The default value is 1.
	Page  string `query:"page"`
	Limit string `query:"limit"`
}

func (q *UserQueryFilter) IsNotEmpty() bool {
	if q.Name == "" && q.Email == "" && q.Phone == "" && q.Password == "" {
		return false
	}

	return true
}