package modelRepository

import "github.com/google/uuid"

type ParamUpsertReview struct {
	AccountUUID  uuid.UUID
	MovieUUID    uuid.UUID
	ReviewRating uint
}

type ParamDeleteReview struct {
	AccountUUID uuid.UUID
	MovieUUID   uuid.UUID
}
