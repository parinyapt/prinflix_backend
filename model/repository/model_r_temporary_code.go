package modelRepository

import (
	"github.com/google/uuid"
	modelDatabase "github.com/parinyapt/prinflix_backend/model/database"
)

type ParamCreateTemporaryCode struct {
	UUID        uuid.UUID
	AccountUUID uuid.UUID
	Type        string
}

type ResultFetchOneTemporaryCode struct {
	IsFound bool
	Data    *modelDatabase.TemporaryCode
}

type ParamDeleteTemporaryCodeByAccountUUIDAndType struct {
	AccountUUID uuid.UUID
	Type        string
}
