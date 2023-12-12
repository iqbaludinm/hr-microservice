package repository

import (
	"context"

	"github.com/iqbaludinm/hr-microservice/profile-service/model/domain"
	"github.com/iqbaludinm/hr-microservice/profile-service/repository/query"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProfileRepository interface {
	// More Priority +1
	CreateUser(c context.Context, user domain.User) error
	UpdateMyProfileTx(ctx context.Context, id string, user domain.User) (domain.User, error)
	FindUserNotDeleteByQueryTx(ctx context.Context, query, value string) (domain.User, error) // bisa dipake forgot-pass
	CheckTokenWithQueryTx(ctx context.Context, query, value string) (domain.ResetPasswordToken, error)
	AddTokenTx(ctx context.Context, token domain.ResetPasswordToken) error
	UpdateTokenTx(ctx context.Context, token domain.ResetPasswordToken) error
	UpdatePasswordTx(ctx context.Context, user domain.User) error

	// -1
	// FindUserWithNameNotDeleteByQueryTx(ctx context.Context, query, value string) (domain.UserWithName, error)
	// FindProjectByIdTx(ctx context.Context, id int) (domain.Project, error)
	// FindDepartementByIdTx(ctx context.Context, id int) (domain.Departement, error)
	// FindRoleByIdTx(ctx context.Context, id int) (domain.Role, error)
	// FindPermissionRoleByRoleIdTx(ctx context.Context, id int) ([]helper.PermissionRole, error)
}

type profileRepository struct {
	db           Store
	ProfileQuery query.ProfileQuery
}

func NewProfile(db Store, q query.ProfileQuery) ProfileRepository {
	return &profileRepository{
		db:           db,
		ProfileQuery: q,
	}
}


func (r *profileRepository) CreateUser(c context.Context, user domain.User) error {
	var err error

	// create transaction to create user
	err = r.db.WithTransaction(c, func(tx pgx.Tx) error {
		// create user, if error will rollback
		if err = r.ProfileQuery.CreateUser(c, tx, user); err != nil {
			return err
		}
		return nil
	})

	return err
}

func (r *profileRepository) FindUserNotDeleteByQueryTx(ctx context.Context, query, value string) (domain.User, error) {

	var data domain.User
	var err error

	err = r.db.WithoutTransaction(ctx, func(db *pgxpool.Pool) error {
		data, err = r.ProfileQuery.FindUserNotDeleteByQuery(ctx, db, query, value)
		if err != nil {
			return err
		}

		return nil
	})

	return data, err
}

func (r *profileRepository) UpdateMyProfileTx(ctx context.Context, id string, user domain.User) (domain.User, error) {
	var data domain.User
	var err error

	err = r.db.WithTransaction(ctx, func(tx pgx.Tx) error {

		data, err = r.ProfileQuery.UpdateMyProfile(ctx, tx, id, user)
		if err != nil {
			return err
		}

		return nil
	})

	return data, err
}

// Token
func (r *profileRepository) CheckTokenWithQueryTx(ctx context.Context, query, value string) (domain.ResetPasswordToken, error) {

	var data domain.ResetPasswordToken
	var err error

	err = r.db.WithTransaction(ctx, func(tx pgx.Tx) error {
		data, err = r.ProfileQuery.CheckTokenWithQuery(ctx, tx, query, value)
		if err != nil {
			return err
		}

		return nil
	})

	return data, err
}

func (r *profileRepository) AddTokenTx(ctx context.Context, token domain.ResetPasswordToken) error {

	var err error
	err = r.db.WithTransaction(ctx, func(tx pgx.Tx) error {

		err = r.ProfileQuery.AddToken(ctx, tx, token)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (r *profileRepository) UpdateTokenTx(ctx context.Context, token domain.ResetPasswordToken) error {

	var err error

	err = r.db.WithTransaction(ctx, func(tx pgx.Tx) error {

		err = r.ProfileQuery.UpdateToken(ctx, tx, token)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (r *profileRepository) UpdatePasswordTx(ctx context.Context, user domain.User) error {

	var err error

	err = r.db.WithTransaction(ctx, func(tx pgx.Tx) error {

		err = r.ProfileQuery.UpdatePassword(ctx, tx, user)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}