package dto

type SignUpReq struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	FullName    string `json:"full_name" binding:"required"`
	GroupNumber string `json:"group_number" binding:"required"`
}
