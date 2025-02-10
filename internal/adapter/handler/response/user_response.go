package response

type UserResponse struct {
	ID    int64  `json:"id"`
	Nama  string `json:"nama"`
	Email string `json:"email"`
}
