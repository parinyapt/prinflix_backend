package repository

import "gorm.io/gorm"

type RepositoryReceiverArgument struct {
	databaseTX *gorm.DB
}

func NewRepository(dbtx *gorm.DB) *RepositoryReceiverArgument {
	return &RepositoryReceiverArgument{
		databaseTX: dbtx,
	}
}

const (
	errorDatabaseQueryFailed = "Database query failed"
)