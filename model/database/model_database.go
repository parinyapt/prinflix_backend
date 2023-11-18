package modelDatabase

import (
	"time"

	"github.com/google/uuid"

	utilsDatabase "github.com/parinyapt/prinflix_backend/utils/database"
)

type Account struct {
	UUID          uuid.UUID `gorm:"column:account_uuid;primary_key;not null"`
	Name          string    `gorm:"column:account_name;not null"`
	Email         string    `gorm:"column:account_email;unique;not null"`
	EmailVerified bool      `gorm:"column:account_email_verified;default:false;not null"`
	Password      string    `gorm:"column:account_password;not null"`
	Status        string    `gorm:"column:account_status;type:enum('active', 'inactive');default:active;not null"`
	Image         bool      `gorm:"column:account_image;default:false;not null"`
	Role          string    `gorm:"column:account_role;type:enum('admin', 'user');default:user;not null"`
	CreatedAt     time.Time `gorm:"column:account_created_at;not null"`
	UpdatedAt     time.Time `gorm:"column:account_updated_at;not null"`
}

func (Account) TableName() string {
	return utilsDatabase.GenerateTableName("account")
}

type AuthSession struct {
	UUID        uuid.UUID `gorm:"column:auth_session_uuid;primary_key;not null"`
	AccountUUID uuid.UUID `gorm:"column:account_uuid;not null"`
	ExpiredAt   time.Time `gorm:"column:auth_session_expired_at;not null"`
	CreatedAt   time.Time `gorm:"column:auth_session_created_at;not null"`
}

func (AuthSession) TableName() string {
	return utilsDatabase.GenerateTableName("auth_session")
}

type AccountOAuth struct {
	AccountUUID uuid.UUID `gorm:"column:account_oauth_account_uuid;primary_key;not null"`
	Provider    string    `gorm:"column:account_oauth_provider;primary_key;type:enum('line', 'google');not null"`
	UserID      string    `gorm:"column:account_oauth_user_id;not null"`
	UserEmail   string    `gorm:"column:account_oauth_user_email;not null"`
	UserPicture string    `gorm:"column:account_oauth_user_picture;not null"`
	CreatedAt   time.Time `gorm:"column:account_oauth_created_at;not null"`
}

func (AccountOAuth) TableName() string {
	return utilsDatabase.GenerateTableName("account_oauth")
}

type TemporaryCode struct {
	UUID        uuid.UUID `gorm:"column:auth_temporary_code_uuid;type:uuid;primary_key;not null"`
	AccountUUID uuid.UUID `gorm:"column:auth_temporary_code_account_uuid;not null"`
	Type        string    `gorm:"column:auth_temporary_code_type;type:enum('email_verification', 'password_reset', 'oauth_state');not null"`
	CreatedAt   time.Time `gorm:"column:auth_temporary_code_created_at;type:timestamp;not null"`
}

func (TemporaryCode) TableName() string {
	return utilsDatabase.GenerateTableName("temporary_code")
}

type Movie struct {
	UUID        uuid.UUID `gorm:"column:movie_uuid;primary_key;not null"`
	CategoryID  uint      `gorm:"column:movie_category_id;not null"`
	Title       string    `gorm:"column:movie_title;not null"`
	Description string    `gorm:"column:movie_description;not null"`
	CreatedAt   time.Time `gorm:"column:movie_created_at;not null"`
	UpdatedAt   time.Time `gorm:"column:movie_updated_at;not null"`
}

func (Movie) TableName() string {
	return utilsDatabase.GenerateTableName("movie")
}

type MovieCategory struct {
	ID        uint      `gorm:"column:movie_category_id;primary_key;not null"`
	Name      string    `gorm:"column:movie_category;not null"`
	CreatedAt time.Time `gorm:"column:movie_category_created_at;not null"`
	UpdatedAt time.Time `gorm:"column:movie_category_updated_at;not null"`
}

func (MovieCategory) TableName() string {
	return utilsDatabase.GenerateTableName("movie_category")
}

type FavoriteMovie struct {
	AccountUUID uuid.UUID `gorm:"column:favorite_movie_account_uuid;primary_key;not null"`
	MovieUUID   uuid.UUID `gorm:"column:favorite_movie_movie_uuid;primary_key;not null"`
	CreatedAt   time.Time `gorm:"column:favorite_movie_created_at;not null"`
}

func (FavoriteMovie) TableName() string {
	return utilsDatabase.GenerateTableName("favorite_movie")
}
