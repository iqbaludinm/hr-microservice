package repository

import (
	"context"

	"github.com/iqbaludinm/hr-microservice/user-service/model/domain"
	"github.com/iqbaludinm/hr-microservice/user-service/repository/query"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	CreateUser(c context.Context, user domain.User) error
	UpdateUser(c context.Context, id string, user domain.User) error
	UpdatePassword(c context.Context, user domain.User) error
	Delete(c context.Context, id string) error
	FindAllUser(c context.Context, filter domain.UserQueryFilter) ([]domain.User, error)
	FindById(c context.Context, id string, filter domain.UserQueryFilter) (domain.User, error)
	FindByPhoneNumber(c context.Context, phone string) (domain.User, error)
	FindByEmail(c context.Context, email string) (domain.User, error)
	FindUserNotDeleteByQueryTx(ctx context.Context, query, value string) (domain.User, error)
	CheckTokenWithQueryTx(ctx context.Context, query, value string) (domain.ResetPasswordToken, error)
	CountAllUser(c context.Context, filter domain.UserQueryFilter) (int, error)
	AddTokenTx(ctx context.Context, token domain.ResetPasswordToken) error
	UpdateTokenTx(ctx context.Context, token domain.ResetPasswordToken) error
	UpdatePasswordTx(ctx context.Context, user domain.User) error
}

type userRepository struct {
	db        Store
	UserQuery query.UserQuery
}

func NewUser(db Store, q query.UserQuery) UserRepository {
	return &userRepository{
		db:        db,
		UserQuery: q,
	}
}

func (r *userRepository) UpdateUser(c context.Context, id string, user domain.User) error {
	var err error

	// create transaction to update user
	err = r.db.WithTransaction(c, func(tx pgx.Tx) error {
		// update user id, if error will rollback
		if err = r.UserQuery.UpdateUser(c, tx, id, user); err != nil {
			return err
		}
		return nil
	})

	return err
}

func (r *userRepository) UpdatePassword(c context.Context, user domain.User) error {
	var err error

	// create transaction to update password user
	err = r.db.WithTransaction(c, func(tx pgx.Tx) error {
		// update password user id, if error will rollback
		if err = r.UserQuery.UpdatePassword(c, tx, user); err != nil {
			return err
		}
		return nil
	})

	return err
}

func (r *userRepository) FindById(c context.Context, id string, filter domain.UserQueryFilter) (domain.User, error) {
	var user domain.User
	var err error

	// get user by Id without transaction.
	err = r.db.WithoutTransaction(c, func(db *pgxpool.Pool) error {
		if user, err = r.UserQuery.FindById(c, db, id, filter); err != nil {
			return err
		}
		return nil
	})

	return user, err
}

// Find User by email
func (r *userRepository) FindByEmail(c context.Context, email string) (domain.User, error) {
	var user domain.User
	var err error

	// get user by email without transaction.
	err = r.db.WithoutTransaction(c, func(db *pgxpool.Pool) error {
		if user, err = r.UserQuery.FindByEmail(c, db, email); err != nil {
			return err
		}
		return nil
	})

	return user, err
}

// Find User by phone number
func (r *userRepository) FindByPhoneNumber(c context.Context, phone string) (domain.User, error) {
	var user domain.User
	var err error

	// get user by phone number without transaction.
	err = r.db.WithoutTransaction(c, func(db *pgxpool.Pool) error {
		if user, err = r.UserQuery.FindByPhoneNumber(c, db, phone); err != nil {
			return err
		}
		return nil
	})

	return user, err
}

func (r *userRepository) FindAllUser(c context.Context, filter domain.UserQueryFilter) ([]domain.User, error) {
	var users []domain.User
	var err error

	// get users by roleID without transaction.
	err = r.db.WithoutTransaction(c, func(db *pgxpool.Pool) error {
		if users, err = r.UserQuery.FindAllUser(c, db, filter); err != nil {
			return err
		}
		return nil
	})

	return users, err
}

func (r *userRepository) CreateUser(c context.Context, user domain.User) error {
	var err error

	// create transaction to create user
	err = r.db.WithTransaction(c, func(tx pgx.Tx) error {
		// create user, if error will rollback
		if err = r.UserQuery.CreateUser(c, tx, user); err != nil {
			return err
		}
		return nil
	})

	return err
}

func (r *userRepository) Delete(c context.Context, id string) error {
	var err error

	// create transaction to delete user
	err = r.db.WithTransaction(c, func(tx pgx.Tx) error {
		// delete user by id, if error will rollback
		if err = r.UserQuery.Delete(c, tx, id); err != nil {
			return err
		}
		return nil
	})

	return err
}

func (r *userRepository) CountAllUser(c context.Context, filter domain.UserQueryFilter) (int, error) {
	var count int
	var err error
	// get Status by ID without transaction.
	err = r.db.WithoutTransaction(c, func(db *pgxpool.Pool) error {
		if count, err = r.UserQuery.CountAllUser(c, db, filter); err != nil {
			return err
		}
		return nil
	})

	return count, err
}


func (r *userRepository) FindUserNotDeleteByQueryTx(ctx context.Context, query, value string) (domain.User, error) {

	var data domain.User
	var err error

	err = r.db.WithoutTransaction(ctx, func(db *pgxpool.Pool) error {
		data, err = r.UserQuery.FindUserNotDeleteByQuery(ctx, db, query, value)
		if err != nil {
			return err
		}

		return nil
	})

	return data, err
}

// Token
func (r *userRepository) CheckTokenWithQueryTx(ctx context.Context, query, value string) (domain.ResetPasswordToken, error) {

	var data domain.ResetPasswordToken
	var err error

	err = r.db.WithTransaction(ctx, func(tx pgx.Tx) error {
		data, err = r.UserQuery.CheckTokenWithQuery(ctx, tx, query, value)
		if err != nil {
			return err
		}

		return nil
	})

	return data, err
}

func (r *userRepository) AddTokenTx(ctx context.Context, token domain.ResetPasswordToken) error {

	var err error
	err = r.db.WithTransaction(ctx, func(tx pgx.Tx) error {

		err = r.UserQuery.AddToken(ctx, tx, token)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (r *userRepository) UpdateTokenTx(ctx context.Context, token domain.ResetPasswordToken) error {

	var err error

	err = r.db.WithTransaction(ctx, func(tx pgx.Tx) error {

		err = r.UserQuery.UpdateToken(ctx, tx, token)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (r *userRepository) UpdatePasswordTx(ctx context.Context, user domain.User) error {

	var err error

	err = r.db.WithTransaction(ctx, func(tx pgx.Tx) error {

		err = r.UserQuery.UpdatePassword(ctx, tx, user)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}