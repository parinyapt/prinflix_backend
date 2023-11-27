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
	SortOrderByAsc  = "ASC"
	SortOrderByDesc = "DESC"
)

const (
	errorDatabaseQueryFailed = "Database query failed"
)
