package dto

type TokensReq struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}
