package request

type ContentRequest struct {
	Title       string `json:"title" validate:"required"`
	Excerp      string `json:"excerp" validate:"required"`
	Description string `json:"description" validate:"required"`
	Image       string `json:"image" validate:"required"`
	Tags        string `json:"tags"`
	CategoryID  int64  `json:"category_id" validate:"required"`
	Status      string `json:"status" validate:"required"`
}
