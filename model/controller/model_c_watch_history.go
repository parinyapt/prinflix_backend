package modelController

type ParamUpdateWatchHistory struct {
	AccountUUID string
	MovieUUID   string
	TimeStamp   int64
}

type ReturnCheckWatchHistory struct {
	IsNotFound bool

	LatestTimeStamp int64
	IsEnd           bool
}
