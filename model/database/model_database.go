package modelDatabase

import (
	"time"

	"github.com/google/uuid"

	utilsDatabase "github.com/parinyapt/prinflix_backend/utils/database"
)

type Account struct {
	UUID      uuid.UUID `gorm:"column:account_uuid;primary_key"`
	Name      string    `gorm:"column:account_name"`
	Image     string    `gorm:"column:account_image"`
	Email     string    `gorm:"column:account_email"`
	Password  string    `gorm:"column:account_password"`
	Status    string    `gorm:"column:account_status;type:enum('active', 'inactive');default:active"`
	CreatedAt time.Time `gorm:"column:account_created_at"`
	UpdatedAt time.Time `gorm:"column:account_updated_at"`
}

func (Account) TableName() string {
	return utilsDatabase.GenerateTableName("account")
}
