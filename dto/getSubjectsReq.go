package dto

type GetSubjectsReq struct {
	Number string `uri:"number" binding:"required"`
}
