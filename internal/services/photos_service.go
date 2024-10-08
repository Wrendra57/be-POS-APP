package services

import (
	"github.com/Wrendra57/Pos-app-be/config"
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/Wrendra57/Pos-app-be/internal/models/webrequest"
	"github.com/Wrendra57/Pos-app-be/internal/repositories"
	"github.com/Wrendra57/Pos-app-be/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type PhotosService interface {
	UploadPhotoService(ctx *fiber.Ctx, request webrequest.PhotoUploadRequest) domain.Photos
	UploadPhoto(ctx *fiber.Ctx, tx pgx.Tx, name string, photo *multipart.FileHeader, owner uuid.UUID) (domain.Photos,
		error)
}

type photosServiceImpl struct {
	PhotoRepo repositories.PhotosRepository
	DB        *pgxpool.Pool
	Validate  *validator.Validate
}

func NewPhotosService(PhotoRepo repositories.PhotosRepository,
	DB *pgxpool.Pool,
	Validate *validator.Validate) PhotosService {
	return &photosServiceImpl{
		PhotoRepo: PhotoRepo,
		DB:        DB,
		Validate:  Validate,
	}
}
func (s photosServiceImpl) UploadPhotoService(ctx *fiber.Ctx, request webrequest.PhotoUploadRequest) domain.Photos {
	request.Name = strings.Join(strings.Split(request.Name, " "), "-")
	request.Foto.Filename = strings.Join(strings.Split(request.Foto.Filename, " "), "-")

	filename := request.Name + "-" + time.Now().Format(
		"20060102_150405") + "-" + request.Foto.Filename

	f := domain.Photos{
		Url:   "http://127.0.0.1:8080/foto/" + filename,
		Owner: request.Owner_id,
	}

	tx, err := s.DB.BeginTx(ctx.Context(), config.TxConfig())
	utils.PanicIfError(err)
	defer utils.CommitOrRollback(ctx.Context(), tx)

	// Buat direktori uploads jika belum ada
	uploadsDir := "./storage/photos"
	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		err := os.MkdirAll(uploadsDir, os.ModePerm)
		utils.PanicIfError(err)
	}

	filepath := filepath.Join(uploadsDir, filename)
	err = ctx.SaveFile(request.Foto, filepath)
	utils.PanicIfError(err)

	f = s.PhotoRepo.Insert(ctx.Context(), tx, f)
	utils.PanicIfError(err)

	return f
}
func (s photosServiceImpl) UploadPhoto(ctx *fiber.Ctx, tx pgx.Tx, name string,
	photo *multipart.FileHeader, owner uuid.UUID) (domain.Photos, error) {
	//TODO implement me
	name = strings.Join(strings.Split(name, " "), "-")
	photo.Filename = strings.Join(strings.Split(photo.Filename, " "), "-")

	filename := name + "-" + time.Now().Format(
		"20060102_150405") + "-" + photo.Filename

	f := domain.Photos{
		Url:   "http://127.0.0.1:8080/foto/" + filename,
		Owner: owner,
	}

	uploadsDir := "./storage/photos"
	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		err := os.MkdirAll(uploadsDir, os.ModePerm)
		utils.PanicIfError(err)
	}
	filepath := filepath.Join(uploadsDir, filename)
	err := ctx.SaveFile(photo, filepath)
	utils.PanicIfError(err)

	f = s.PhotoRepo.Insert(ctx.Context(), tx, f)
	utils.PanicIfError(err)

	return f, nil
}
