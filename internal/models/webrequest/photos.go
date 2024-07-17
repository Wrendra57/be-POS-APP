package webrequest

import (
	"github.com/google/uuid"
	"mime/multipart"
)

type PhotoUploadRequest struct {
	Foto     *multipart.FileHeader `json:"foto" validate:"required"`
	Owner_id uuid.UUID             `json:"owner_id" validate:"required"`
	Name     string                `json:"name" validate:"required"`
}
