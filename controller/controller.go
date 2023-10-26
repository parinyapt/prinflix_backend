package controller

import "gorm.io/gorm"

type ControllerReceiverArgument struct {
	databaseTX *gorm.DB
}

func NewController(dbtx *gorm.DB) *ControllerReceiverArgument {
	return &ControllerReceiverArgument{
		databaseTX: dbtx,
	}
}