package modelDatabase

import (
	"time"

	"github.com/google/uuid"

	utilsDatabase "github.com/parinyapt/prinflix_backend/utils/database"
)

const (
	AccountStatusActive   = "active"
	AccountStatusInactive = "inactive"

	AccountRoleAdmin = "admin"
	AccountRoleUser  = "user"

	AccountOAuthProviderLine   = "line"
	AccountOAuthProviderGoogle = "google"
	AccountOAuthProviderApple  = "apple"

	TemporaryCodeTypeEmailVerification = "email_verification"
	TemporaryCodeTypePasswordReset     = "password_reset"
	TemporaryCodeTypeOAuthStateLine    = "oauth_state_line"
	TemporaryCodeTypeOAuthStateGoogle  = "oauth_state_google"
	TemporaryCodeTypeOAuthStateApple   = "oauth_state_apple"
	TemporaryCodeTypeAuthTokenCode     = "auth_token_code"

	OauthStateProviderLine   = "line"
	OauthStateProviderGoogle = "google"
	OauthStateProviderApple  = "apple"

	ReviewRatingGood = 3
	ReviewRatingFair = 2
	ReviewRatingBad  = 1
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
	Provider    string    `gorm:"column:account_oauth_provider;primary_key;type:enum('line', 'google', 'apple');not null"`
	UserID      string    `gorm:"column:account_oauth_user_id;not null"`
	UserName    string    `gorm:"column:account_oauth_user_name;not null"`
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
	Type        string    `gorm:"column:auth_temporary_code_type;type:enum('email_verification', 'password_reset', 'oauth_state_line', 'oauth_state_google', 'oauth_state_apple', 'auth_token_code');not null"`
	CreatedAt   time.Time `gorm:"column:auth_temporary_code_created_at;type:timestamp;not null"`
}

func (TemporaryCode) TableName() string {
	return utilsDatabase.GenerateTableName("temporary_code")
}

type Movie struct {
	UUID        uuid.UUID `gorm:"column:movie_uuid;primary_key;not null"`
	CategoryID  uint      `gorm:"column:movie_movie_category_id;not null"`
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

type WatchSession struct {
	UUID        uuid.UUID `gorm:"column:watch_session_uuid;primary_key;not null"`
	AccountUUID uuid.UUID `gorm:"column:watch_session_account_uuid;not null"`
	MovieUUID   uuid.UUID `gorm:"column:watch_session_movie_uuid;not null"`
	ExpiredAt   time.Time `gorm:"column:watch_session_expired_at;not null"`
	CreatedAt   time.Time `gorm:"column:watch_session_created_at;not null"`
}

func (WatchSession) TableName() string {
	return utilsDatabase.GenerateTableName("watch_session")
}

type WatchHistory struct {
	UUID            uuid.UUID `gorm:"column:watch_history_uuid;primary_key;not null"`
	AccountUUID     uuid.UUID `gorm:"column:watch_history_account_uuid;not null"`
	MovieUUID       uuid.UUID `gorm:"column:watch_history_movie_uuid;not null"`
	LatestTimeStamp int64     `gorm:"column:watch_history_latest_timestamp;default:0;not null"`
	IsEnd           bool      `gorm:"column:watch_history_is_end;default:false;not null"`
	CreatedAt       time.Time `gorm:"column:watch_history_created_at;not null"`
	UpdatedAt       time.Time `gorm:"column:watch_history_updated_at;not null"`
}

func (WatchHistory) TableName() string {
	return utilsDatabase.GenerateTableName("watch_history")
}

type OauthState struct {
	UUID      uuid.UUID `gorm:"column:oauth_state_uuid;primary_key;not null"`
	Provider  string    `gorm:"column:oauth_state_provider;type:enum('line', 'google', 'apple');not null"`
	CreatedAt time.Time `gorm:"column:oauth_state_created_at;not null"`
}

func (OauthState) TableName() string {
	return utilsDatabase.GenerateTableName("oauth_state")
}

type Review struct {
	AccountUUID uuid.UUID `gorm:"column:review_account_uuid;primary_key;not null"`
	MovieUUID   uuid.UUID `gorm:"column:review_movie_uuid;primary_key;not null"`
	Rating      uint      `gorm:"column:review_rating;not null"`
	CreatedAt   time.Time `gorm:"column:review_created_at;not null"`
	UpdatedAt   time.Time `gorm:"column:review_updated_at;not null"`
}

func (Review) TableName() string {
	return utilsDatabase.GenerateTableName("review")
}
