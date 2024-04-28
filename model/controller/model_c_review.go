package modelController

type ParamCreateUpdateReviewMovie struct {
	AccountUUID  string
	MovieUUID    string
	ReviewRating uint
}

type ParamDeleteReviewMovie struct {
	AccountUUID string
	MovieUUID   string
}
