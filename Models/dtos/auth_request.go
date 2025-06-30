package dtos

type AuthRequest struct {
	Email    string `json:"email"`
	Nickname string `json:"nick"`
	Passwd   string `json:"password"`
}
