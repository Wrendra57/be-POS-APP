package services

import (
	"fmt"
	"github.com/Wrendra57/Pos-app-be/config"
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/Wrendra57/Pos-app-be/internal/models/webrequest"
	"github.com/Wrendra57/Pos-app-be/internal/repositories"
	"github.com/Wrendra57/Pos-app-be/internal/utils"
	"github.com/Wrendra57/Pos-app-be/internal/utils/exception"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type PhotosService interface {
	UploadPhotos(ctx *fiber.Ctx, request webrequest.PhotoUploadRequest) (domain.Photos, exception.CustomEror, bool)
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
func (s photosServiceImpl) UploadPhotos(ctx *fiber.Ctx, request webrequest.PhotoUploadRequest) (domain.Photos,
	exception.CustomEror, bool) {
	request.Name = strings.Join(strings.Split(request.Name, " "), "-")
	request.Foto.Filename = strings.Join(strings.Split(request.Foto.Filename, " "), "-")

	filename := request.Name + "-" + time.Now().Format(
		"20060102_150405") + "-" + request.Foto.Filename

	f := domain.Photos{
		Url:       "http://127.0.0.1:8080/foto/" + filename,
		Owner:     request.Owner_id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tx, err := s.DB.BeginTx(ctx.Context(), config.TxConfig())
	utils.PanicIfError(err)
	defer utils.CommitOrRollback(ctx, tx)

	// Buat direktori uploads jika belum ada
	uploadsDir := "./storage/photos"
	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		err := os.Mkdir(uploadsDir, os.ModePerm)
		if err != nil {
			return domain.Photos{}, exception.CustomEror{Code: 500, Error: err.Error()}, false
		}
	}

	filepath := filepath.Join(uploadsDir, filename)
	err = ctx.SaveFile(request.Foto, filepath)
	utils.PanicIfError(err)

	fmt.Println(filepath)

	f, err = s.PhotoRepo.Insert(ctx, tx, f)
	utils.PanicIfError(err)

	return f, exception.CustomEror{}, true
}
