package dtos

type UpdateAccountResponse struct {
	ID    int    `json:"id" binding:"required"`
	Email string `json:"email" binding:"omitempty,email"`
}
