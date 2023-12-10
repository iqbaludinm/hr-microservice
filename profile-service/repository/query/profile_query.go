package query

import (
	"context"
	"fmt"

	"github.com/iqbaludinm/hr-microservice/profile-service/model/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProfileQuery interface {
	// Register(c context.Context, tx pgx.Tx, user domain.User) (string, error)
	// Login(c context.Context, tx pgx.Tx, id string) (domain.User, error)
	CreateUser(c context.Context, tx pgx.Tx, user domain.User) error
	UpdateMyProfile(c context.Context, tx pgx.Tx, id string, user domain.User) (domain.User, error)
	FindUserNotDeleteByQuery(c context.Context, db *pgxpool.Pool, query, value string) (domain.User, error) // forgot-pass
	// CheckTokenWithQuery(ctx context.Context, db pgx.Tx, query, value string) (domain.ResetPasswordToken, error)
	// AddToken(ctx context.Context, db pgx.Tx, tokens domain.ResetPasswordToken) error
	// UpdateToken(ctx context.Context, db pgx.Tx, tokens domain.ResetPasswordToken) error
}

type ProfileQueryImpl struct {
}

func NewProfile() ProfileQuery {
	return &ProfileQueryImpl{}
}

func (repository *ProfileQueryImpl) CreateUser(c context.Context, tx pgx.Tx, user domain.User) error {
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

func (repository *ProfileQueryImpl) Login(c context.Context, tx pgx.Tx, email string) (domain.User, error) {
	queryStr := fmt.Sprintf("SELECT * FROM %s WHERE email = $1 AND deleted_at is NULL", "users")

	user, err := tx.Query(context.Background(), queryStr, email)
	if err != nil {
		return domain.User{}, err
	}

	defer user.Close()
	data, err := pgx.CollectOneRow(user, pgx.RowToStructByPos[domain.User])

	if err != nil {
		return domain.User{}, err
	}

	return data, nil
}

func (repository *ProfileQueryImpl) FindUserNotDeleteByQuery(c context.Context, db *pgxpool.Pool, query, value string) (domain.User, error) {
	var data domain.User
	queryStr := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1 AND deleted_at is NULL", "users", query)

	row := db.QueryRow(c, queryStr, value)

	err := row.Scan(&data.ID, &data.Name, &data.Email, &data.Password, &data.Phone, &data.CreatedAt, &data.UpdatedAt, &data.DeletedAt)

	if err != nil {
		return domain.User{}, err
	}

	return data, nil
}

func (repository *ProfileQueryImpl) UpdateMyProfile(c context.Context, tx pgx.Tx, id string, user domain.User) (domain.User, error) {
	var data domain.User
	query := "UPDATE users SET name = $1, email = $2, phone = $3 WHERE id = $4"
	_, err := tx.Prepare(context.Background(), "update_profile", query)
	if err != nil {
		return domain.User{}, err
	}

	_, err = tx.Exec(context.Background(), "update_profile", user.Name, user.Email, user.Phone, id)
	if err != nil {
		return domain.User{}, err
	}

	return data, nil
}

// func (repository *ProfileQueryImpl) FindUserWithNameNotDeleteByQuery(ctx context.Context, db pgx.Tx, query, value string) (domain.UserWithName, error) {
// 	queryStr := fmt.Sprintf(`SELECT u.*, coalesce(r.name,''), coalesce(d.name,''), coalesce(p.name,''), coalesce(pst.name,'') FROM %s u
// 	LEFT JOIN roles r ON u.role_id = r.id
// 	LEFT JOIN departements d ON u.departement_id = d.id
// 	LEFT JOIN projects p ON u.project_id = p.id
// 	LEFT JOIN positions pst ON pst.id = u.position_id
// 	WHERE %s = $1 AND deleted_at is NULL`, "users", query)

// 	user, err := db.Query(context.Background(), queryStr, value)

// 	if err != nil {
// 		return domain.UserWithName{}, err
// 	}

// 	defer user.Close()

// 	data, err := pgx.CollectOneRow(user, pgx.RowToStructByPos[domain.UserWithName])

// 	if err != nil {
// 		return domain.UserWithName{}, err
// 	}

// 	return data, nil
// }

// Token
// func (repository *ProfileQueryImpl) CheckTokenWithQuery(ctx context.Context, db pgx.Tx, query, value string) (domain.ResetPasswordToken, error) {
// 	queryStr := fmt.Sprintf("SELECT * FROM %s WHERE %s = '%s'", "reset_token", query, value)

// 	user := db.QueryRow(context.Background(), queryStr)

// 	var data domain.ResetPasswordToken
// 	err := user.Scan(&data.Id, &data.Tokens, &data.Email, &data.Attempt, &data.LastAttempt)
// 	if err != nil {
// 		return domain.ResetPasswordToken{}, err
// 	}

// 	return data, nil
// }

// func (repository *ProfileQueryImpl) AddToken(ctx context.Context, db pgx.Tx, tokens domain.ResetPasswordToken) error {
// 	query := fmt.Sprintf("INSERT INTO %s (id, tokens, email, attempt, last_attempt) VALUES($1,$2,$3,$4, $5)", "reset_token")
// 	_, err := db.Prepare(context.Background(), "add_token", query)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = db.Exec(context.Background(), "add_token", tokens.Id, tokens.Tokens, tokens.Email, tokens.Attempt, tokens.LastAttempt)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (repository *ProfileQueryImpl) UpdateToken(ctx context.Context, db pgx.Tx, tokens domain.ResetPasswordToken) error {
// 	query := fmt.Sprintf("UPDATE %s SET tokens = $1, attempt = $2, last_attempt = $3 WHERE email = $4", "reset_token")
// 	_, err := db.Prepare(context.Background(), "update_token", query)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = db.Exec(context.Background(), "update_token", tokens.Tokens, tokens.Attempt, tokens.LastAttempt, tokens.Email)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
