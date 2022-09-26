package dto

import (
	"github.com/lib/pq"
	"mime/multipart"
)

type UpdateMultipleWithFileReq struct {
	Ids  pq.StringArray        `json:"ids" binding:"required"`
	File *multipart.FileHeader `form:"file" binding:"required"`
}
