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
	CreatedAt     time.Time `gorm:"column:account_created_at;not null"`
	UpdatedAt     time.Time `gorm:"column:account_updated_at;not null"`
}

func (Account) TableName() string {
	return utilsDatabase.GenerateTableName("account")
}
