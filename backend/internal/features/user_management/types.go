package usermanagement

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Username  string
	Email     string
	Phone     string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(username, email, phone string) *User {
	return &User{
		ID:        uuid.New(),
		Username:  username,
		Email:     email,
		Phone:     phone,
		IsActive:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

type AuthCredential struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	Provider     string    `json:"provider,omitempty"`
	ProviderUID  string    `json:"provider_uid,omitempty"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at"`
}

type UserRepositoryInterface interface {
	GetAllUsers(ctx context.Context /*pagination, filter*/) ([]User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*User, error)
	GetUserByCredential(ctx context.Context, username string) (*User, error)
	GetUserByProviderUID(ctx context.Context, provider, providerUID string) (*User, error)
	CheckUserCredentials(ctx context.Context, provider, providerUID string, password string) (*User, error)
	CreateUser(ctx context.Context, user User) (uuid.UUID, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, updated User) (uuid.UUID, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
	SoftDeleteUser(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserListResponse struct {
	UserList []UserResponse `json:"users"`
}

type UserRegisterRequest struct {
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Provider    string `json:"provider"`
	ProviderUID string `json:"provider_uid"`
	Password    string `json:"password"`
	Username    string `json:"username"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserServiceInterface interface {
	ListUsers(ctx context.Context /*pagination,filter*/) (UserListResponse, error)
	RegisterUser(ctx context.Context, userRegisterRequest UserRegisterRequest) (string, error)
	Login(ctx context.Context, loginRequest LoginRequest) (accessToken string, refreshToken string, err error)

	//Logout(ctx context.Context, logoutRequest LogoutRequest) error
	GetUserInfo(ctx context.Context, username string, userID uuid.UUID) (UserResponse, error)
}

type UserHandlerInterface interface {
	GetUserInfo(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
}
