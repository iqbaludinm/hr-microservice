package repository

import (
	"context"

	"github.com/iqbaludinm/hr-microservice/auth-service/model/domain"
	"github.com/iqbaludinm/hr-microservice/auth-service/repository/query"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository interface {
	// More Priority +1 
	RegisterTx(ctx context.Context, user domain.User) (string, error)
	LoginTx(ctx context.Context, email string) (domain.User, error)
	UpdatePasswordTx(ctx context.Context, user domain.User) error
	FindUserNotDeleteByQueryTx(ctx context.Context, query, value string) (domain.User, error) // bisa dipake forgot-pass
	CheckTokenWithQueryTx(ctx context.Context, query, value string) (domain.ResetPasswordToken, error)
	AddTokenTx(ctx context.Context, token domain.ResetPasswordToken) error
	UpdateTokenTx(ctx context.Context, token domain.ResetPasswordToken) error
	UpdateUser(c context.Context, id string, user domain.User) error
	// -1
	// FindUserWithNameNotDeleteByQueryTx(ctx context.Context, query, value string) (domain.UserWithName, error)
	// FindProjectByIdTx(ctx context.Context, id int) (domain.Project, error)
	// FindDepartementByIdTx(ctx context.Context, id int) (domain.Departement, error)
	// FindRoleByIdTx(ctx context.Context, id int) (domain.Role, error)
	// FindPermissionRoleByRoleIdTx(ctx context.Context, id int) ([]helper.PermissionRole, error)
}

type authRepository struct {
	db        Store
	AuthQuery query.AuthQuery
}

func NewAuth(db Store, q query.AuthQuery) AuthRepository {
	return &authRepository{
		db:        db,
		AuthQuery: q,
	}
}

func (r *authRepository) RegisterTx(ctx context.Context, user domain.User) (string, error) {
	var id string
	var err error
	err = r.db.WithTransaction(ctx, func(tx pgx.Tx) error {
		id, err = r.AuthQuery.Register(ctx, tx, user)
		if err != nil {
			return err
		}

		return nil
	})

	return id, err
}

func (r *authRepository) LoginTx(ctx context.Context, email string) (domain.User, error) {

	var data domain.User
	var err error

	err = r.db.WithTransaction(ctx, func(tx pgx.Tx) error {
		data, err = r.AuthQuery.Login(ctx, tx, email)
		if err != nil {
			return err
		}

		return nil
	})

	return data, err
}

func (r *authRepository) FindUserNotDeleteByQueryTx(ctx context.Context, query, value string) (domain.User, error) {

	var data domain.User
	var err error

	err = r.db.WithoutTransaction(ctx, func(db *pgxpool.Pool) error {
		data, err = r.AuthQuery.FindUserNotDeleteByQuery(ctx, db, query, value)
		if err != nil {
			return err
		}

		return nil
	})

	return data, err
}

func (r *authRepository) UpdatePasswordTx(ctx context.Context, user domain.User) error {

	var err error

	err = r.db.WithTransaction(ctx, func(tx pgx.Tx) error {

		err = r.AuthQuery.UpdatePassword(ctx, tx, user)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

// func (r *AuthRepositoryImpl) FindUserWithNameNotDeleteByQueryTx(ctx context.Context, query, value string) (domain.UserWithName, error) {

// 	var data domain.UserWithName
// 	var err error

// 	err = repository.DB.WithTransaction(ctx, func(tx pgx.Tx) error {

// 		data, err = repository.FindUserWithNameNotDeleteByQuery(ctx, tx, query, value)
// 		if err != nil {
// 			return err
// 		}

// 		return nil
// 	})

// 	return data, err
// }

// Token
func (r *authRepository) CheckTokenWithQueryTx(ctx context.Context, query, value string) (domain.ResetPasswordToken, error) {

	var data domain.ResetPasswordToken
	var err error

	err = r.db.WithTransaction(ctx, func(tx pgx.Tx) error {
		data, err = r.AuthQuery.CheckTokenWithQuery(ctx, tx, query, value)
		if err != nil {
			return err
		}

		return nil
	})

	return data, err
}

func (r *authRepository) AddTokenTx(ctx context.Context, token domain.ResetPasswordToken) error {
	
	var err error
	err = r.db.WithTransaction(ctx, func(tx pgx.Tx) error {

		err = r.AuthQuery.AddToken(ctx, tx, token)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (r *authRepository) UpdateTokenTx(ctx context.Context, token domain.ResetPasswordToken) error {

	var err error

	err = r.db.WithTransaction(ctx, func(tx pgx.Tx) error {

		err = r.AuthQuery.UpdateToken(ctx, tx, token)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (r *authRepository) UpdateUser(c context.Context, id string, user domain.User) error {
	var err error

	// create transaction to update user
	err = r.db.WithTransaction(c, func(tx pgx.Tx) error {
		// update user id, if error will rollback
		if err = r.AuthQuery.UpdateUser(c, tx, id, user); err != nil {
			return err
		}
		return nil
	})

	return err
}


// // Project
// func (r *AuthRepositoryImpl) FindProjectByIdTx(ctx context.Context, id int) (domain.Project, error) {

// 	var data domain.Project
// 	var err error

// 	err = repository.DB.WithTransaction(ctx, func(tx pgx.Tx) error {

// 		data, err = repository.FindProjectById(ctx, tx, id)
// 		if err != nil {
// 			return err
// 		}

// 		return nil
// 	})

// 	return data, err
// }

// // Departement
// func (r *AuthRepositoryImpl) FindDepartementByIdTx(ctx context.Context, id int) (domain.Departement, error) {

// 	var data domain.Departement
// 	var err error

// 	err = repository.DB.WithTransaction(ctx, func(tx pgx.Tx) error {

// 		data, err = repository.FindDepartementById(ctx, tx, id)
// 		if err != nil {
// 			return err
// 		}

// 		return nil
// 	})

// 	return data, err
// }

// // Role
// func (r *AuthRepositoryImpl) FindRoleByIdTx(ctx context.Context, id int) (domain.Role, error) {

// 	var data domain.Role
// 	var err error

// 	err = repository.DB.WithTransaction(ctx, func(tx pgx.Tx) error {

// 		data, err = repository.FindRoleById(ctx, tx, id)
// 		if err != nil {
// 			return err
// 		}

// 		return nil
// 	})

// 	return data, err
// }

// // Role
// func (r *AuthRepositoryImpl) FindPermissionRoleByRoleIdTx(ctx context.Context, id int) ([]helper.PermissionRole, error) {

// 	var data []helper.PermissionRole
// 	var err error

// 	err = repository.DB.WithTransaction(ctx, func(tx pgx.Tx) error {

// 		data, err = repository.FindRolePermissionByRoleId(ctx, tx, id)
// 		if err != nil {
// 			return err
// 		}

// 		return nil
// 	})

// 	return data, err
// }
