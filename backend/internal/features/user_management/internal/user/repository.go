package user

import (
	"context"
	"database/sql"
	usermanagement "dukkanim-api/internal/features/user_management"

	"github.com/google/uuid"
)

type SqlUserRepository struct {
	db *sql.DB
}

const (
	getAllUserQuery string = `
	SELECT * FROM user_management.users
	`

	getUserByIdQuery string = `
	SELECT * 
	FROM user_management.users
	WHERE user.id=$1
	`

	getUserByCredentialQuery string = `
	SELECT u.*
	FROM user_management.users AS u
	JOIN user_management.auth_credentials AS ac
	ON u.id=ac.user_id
	WHERE ac.provider=$1 AND ac.user_id=$2
	`
	//TODO: burada provider sadece local olabilir gibi geldi

	getUserByProviderUIDQuery string = `
	SELECT u.*
	FROM user_management.users AS u
	JOIN user_management.auth_credentials AS ac
	ON u.id=ac.user_id
	WHERE ac.provider=$1 AND ac.provider_uid=$2
	`

	getLocaleAuthCredentialsQuery string = `
	SELECT ac.id, ac.user_id, ac.password_hash
	FROM user_management.auth_credentials AS ac
	JOIN user_management.users AS u
	ON u.id=ac.user_id
	WHERE u.username=$1 AND ac.provider=$2
	`
	createUserQuery string = ``

	createUserCredentialQuery string = ``

	updateUserQuery string = ``

	softDeleteUserQuery string = ``
)

func (r SqlUserRepository) GetAllUsers(ctx context.Context /*pagination, filter*/) ([]usermanagement.User, error) {

}
func (r SqlUserRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (*usermanagement.User, error) {
}
func (r SqlUserRepository) GetUserByCredential(ctx context.Context, phone, email string) (*usermanagement.User, error) {
}
func (r SqlUserRepository) GetUserByProviderUID(ctx context.Context, provider, providerUID string) (*usermanagement.User, error) {
}
func (r SqlUserRepository) CheckUserCredentials(ctx context.Context, provider, providerUID string, password string) error {
}
func (r SqlUserRepository) CreateUser(ctx context.Context, user usermanagement.User) (uuid.UUID, error) {
}
func (r SqlUserRepository) UpdateUser(ctx context.Context, userID uuid.UUID, updated usermanagement.User) (uuid.UUID, error) {
}
func (r SqlUserRepository) DeleteUser(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {}

func (r SqlUserRepository) SoftDeleteUser(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {}
