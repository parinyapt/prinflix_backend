package modelHandler

type RequestPauseMovie struct {
	TimeStamp int64 `json:"timestamp" validate:"required,min=0,number"`
}
