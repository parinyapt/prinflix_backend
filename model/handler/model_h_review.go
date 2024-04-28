package modelHandler

type RequestReviewMovie struct {
	Rating string `json:"rating" validate:"min=0,max=4"`
}