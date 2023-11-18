package modelController

import "github.com/google/uuid"

type ParamTemporaryCode struct {
	AccountUUID string
	Type        string
}

type ReturnCreateTemporaryCode struct {
	CodeUUID uuid.UUID
}

type ParamCheckTemporaryCode struct {
	CodeUUID string
	Type     string
}

type ReturnCheckTemporaryCode struct {
	IsNotFound  bool
	IsExpired   bool
	AccountUUID uuid.UUID
	CodeUUID    uuid.UUID
}
