package query

import (
	"context"
	"fmt"
	"time"

	"github.com/iqbaludinm/hr-microservice/user-service/model/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserQuery interface {
	CreateUser(c context.Context, tx pgx.Tx, user domain.User) error
	UpdateUser(c context.Context, tx pgx.Tx, id string, user domain.User) error
	UpdatePassword(c context.Context, tx pgx.Tx, user domain.User) error
	Delete(c context.Context, tx pgx.Tx, id string) error
	FindAllUser(c context.Context, db *pgxpool.Pool, filter domain.UserQueryFilter) ([]domain.User, error)
	FindById(c context.Context, db *pgxpool.Pool, id string, filter domain.UserQueryFilter) (domain.User, error)
	FindByEmail(c context.Context, db *pgxpool.Pool, email string) (domain.User, error)
	FindByPhoneNumber(c context.Context, db *pgxpool.Pool, phone string) (domain.User, error)

	CountAllUser(c context.Context, db *pgxpool.Pool, filter domain.UserQueryFilter) (int, error)
}

type UserQueryImpl struct {
}

func NewUser() UserQuery {
	return &UserQueryImpl{}
}

func (repository *UserQueryImpl) UpdateUser(c context.Context, tx pgx.Tx, id string, user domain.User) error {
	// build UPDATE query
	query := fmt.Sprintf("UPDATE %s SET name=$1, email=$2, phone=$3, updated_at=$4 WHERE id=$5", "users")

	_, err := tx.Exec(c, query, user.Name, user.Email, user.Phone, user.UpdatedAt, id)

	return err
}

func (repository *UserQueryImpl) UpdatePassword(c context.Context, tx pgx.Tx, user domain.User) error {
	// build UPDATE query
	query := `UPDATE users SET password=$1 WHERE id=$2`

	_, err := tx.Exec(c, query, user.Password)
	if err != nil {
		return err
	}

	return err
}

func (repository *UserQueryImpl) FindById(c context.Context, db *pgxpool.Pool, id string, filter domain.UserQueryFilter) (domain.User, error) {
	// user query filter builders
	filterString, _ := filter.BuildUserQueries()

	// build SELECT query
	query := fmt.Sprintf(
		`SELECT 
			u.id, 
			u.name, 
			u.email, 
			u.password, 
			u.phone, 
			u.created_at, 
			u.updated_at, 
			u.deleted_at
		FROM users AS u
		WHERE
			%s
			AND id=$1
		`,
		filterString,
	)
	row := db.QueryRow(c, query, id)

	var data domain.User
	err := row.Scan(&data.ID, &data.Name, &data.Email, &data.Password, &data.Phone, &data.CreatedAt, &data.UpdatedAt, &data.DeletedAt)

	return data, err
}

func (repository *UserQueryImpl) FindByEmail(c context.Context, db *pgxpool.Pool, email string) (domain.User, error) {
	// build SELECT query
	query := `
		SELECT u.id, u.name, u.email, u.password, u.phone, u.created_at, u.updated_at, u.deleted_at
			FROM users AS u
			WHERE
				u.deleted_at is null AND
				email=$1
		`
	row := db.QueryRow(c, query, email)

	var data domain.User
	err := row.Scan(&data.ID, &data.Name, &data.Email, &data.Password, &data.Phone, &data.CreatedAt, &data.UpdatedAt, &data.DeletedAt)

	return data, err
}

func (repository *UserQueryImpl) FindByPhoneNumber(c context.Context, db *pgxpool.Pool, phone string) (domain.User, error) {
	// build SELECT query
	query := `
		SELECT u.id, u.name, u.email, u.password, u.phone, u.created_at, u.updated_at, u.deleted_at
			FROM users AS u 
			WHERE
				u.deleted_at is null AND
				phone=$1
		`
	row := db.QueryRow(c, query, phone)

	var data domain.User
	err := row.Scan(&data.ID, &data.Name, &data.Email, &data.Password, &data.Phone, &data.CreatedAt, &data.UpdatedAt, &data.DeletedAt)

	return data, err
}

func (repository *UserQueryImpl) FindAllUser(c context.Context, db *pgxpool.Pool, filter domain.UserQueryFilter) ([]domain.User, error) {
	// user query filter builders
	filterString, _ := filter.BuildUserQueries()

	query := fmt.Sprintf(
		`SELECT * FROM users as u %s`, filterString,
	)
	rows, err := db.Query(c, query)
	if err != nil {
		return []domain.User{}, err
	}
	defer rows.Close()

	var datas []domain.User
	for rows.Next() {
		var data domain.User
		err := rows.Scan(&data.ID, &data.Name, &data.Email, &data.Password, &data.Phone, &data.CreatedAt, &data.UpdatedAt, &data.DeletedAt)
		if err != nil {
			return []domain.User{}, err
		}
		datas = append(datas, data)

	}

	return datas, err
}

// KAFKA INTEGRATION: inserting new data user to database
// UNUSED FUNC, pake modul auth
func (repository *UserQueryImpl) CreateUser(c context.Context, tx pgx.Tx, user domain.User) error {
	// build INSERT query
	query := `INSERT INTO users (
		"id", 
		"name", 
		"email", 
		"password", 
		"phone", 
		"created_at", 
		"updated_at", 
		"deleted_at"
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`

	_, err := tx.Exec(c, query,
		user.ID,
		user.Name,
		user.Email,
		user.Password,
		user.Phone,
		user.CreatedAt,
		user.UpdatedAt,
		user.DeletedAt,
	)

	return err
}

func (repository *UserQueryImpl) Delete(c context.Context, tx pgx.Tx, id string) error {
	// build UPDATE query
	query := `UPDATE users SET deleted_at=$1 WHERE id=$1`

	_, err := tx.Exec(c, query, time.Now(), id)

	return err
}

func (r *UserQueryImpl) CountAllUser(c context.Context, db *pgxpool.Pool, filter domain.UserQueryFilter) (int, error) {
	// user query filter builders
	filterString, _ := filter.BuildUserQueries()

	query := fmt.Sprintf(
		`SELECT
			COUNT(*)
		FROM users AS u
		%s
		`,
		filterString,
	)

	row := db.QueryRow(c, query)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, err
}
