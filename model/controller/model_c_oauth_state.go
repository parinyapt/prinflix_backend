package modelController

import "github.com/google/uuid"

type ReturnCreateOauthState struct {
	StateUUID uuid.UUID
}

type ParamOauthState struct {
	StateUUID string
	Provider  string
}

type ReturnCheckOauthState struct {
	IsNotFound bool
	IsExpired  bool
	StateUUID  uuid.UUID
}
