package dtos

type CreateAccountResponse struct {
	ID    int    `json:"id" binding:"required"`
	Nick  string `json:"nick" binding:"required"`
	Email string `json:"email" binding:"required"`
}
