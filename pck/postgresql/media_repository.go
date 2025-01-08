package postgresql

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/cosmintimis/deepfake-guardian-api/pck/business/models"
	"github.com/cosmintimis/deepfake-guardian-api/pck/business/repositories"
	"github.com/cosmintimis/deepfake-guardian-api/pck/utils"
	"github.com/google/uuid"
)

type mediaRespository struct {
	logger *slog.Logger
}

func NewMediaRepository(logger *slog.Logger) repositories.MediaRepository {

	return &mediaRespository{
		logger: logger,
	}
}

func (mr *mediaRespository) GetByID(id string) (*models.Media, error) {
	dbConnection := GetDBConnection()
	if dbConnection == nil {
		mr.logger.Error("failed to get db connection")
		return nil, fmt.Errorf("failed to get db connection")
	}
	var media models.Media
	err := dbConnection.QueryRow(context.Background(), "SELECT id, title, description, location, type, mimeType, size, tags, mediaData FROM media WHERE id = $1", id).Scan(&media.Id, &media.Title, &media.Description, &media.Location, &media.Type, &media.MimeType, &media.Size, &media.Tags, &media.MediaData)
	if err != nil {
		mr.logger.Error("failed to get media by id", slog.Any("error", err))
		return nil, utils.ErrMediaNotFound
	}
	return &media, nil
}

func (mr *mediaRespository) Create(media *repositories.MediaPayload) (*models.Media, error) {
	dbConnection := GetDBConnection()
	if dbConnection == nil {
		mr.logger.Error("failed to get db connection")
		return nil, fmt.Errorf("failed to get db connection")
	}
	generatedId := uuid.NewString()
	_, err := dbConnection.Exec(context.Background(), "INSERT INTO media (id, title, description, location, type, mimeType, size, tags, mediaData) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)", generatedId, media.Title, media.Description, media.Location, media.Type, media.MimeType, media.Size, media.Tags, media.MediaData)
	if err != nil {
		mr.logger.Error("failed to create media", slog.Any("error", err))
		return nil, fmt.Errorf("failed to create media: %w", err)
	}

	// Retrieve the inserted media
	var createdMedia models.Media
	err = dbConnection.QueryRow(context.Background(),
		"SELECT id, title, description, location, type, mimeType, size, tags, mediaData FROM media WHERE id = $1",
		generatedId).Scan(&createdMedia.Id, &createdMedia.Title, &createdMedia.Description, &createdMedia.Location, &createdMedia.Type, &createdMedia.MimeType, &createdMedia.Size, &createdMedia.Tags, &createdMedia.MediaData)

	if err != nil {
		mr.logger.Error("failed to retrieve created media", slog.Any("error", err))
		return nil, fmt.Errorf("failed to retrieve created media: %w", err)
	}

	return &createdMedia, nil
}

func (mr *mediaRespository) Update(id string, media *repositories.MediaPayload) (*models.Media, error) {
	dbConnection := GetDBConnection()
	if dbConnection == nil {
		mr.logger.Error("failed to get db connection")
		return nil, fmt.Errorf("failed to get db connection")
	}
	// look if the media exists
	var mediaExists bool
	err := dbConnection.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM media WHERE id = $1)", id).Scan(&mediaExists)
	if err != nil {
		mr.logger.Error("failed to check if media exists", slog.Any("error", err))
		return nil, fmt.Errorf("failed to check if media exists: %w", err)
	}
	if !mediaExists {
		mr.logger.Error("media with id does not exist", slog.Any("id", id))
		return nil, utils.ErrMediaNotFound
	}

	// Prepare query parts and values for non-empty fields
	updateFields := []string{}
	values := []interface{}{}

	if media.Title != "" {
		updateFields = append(updateFields, "title = $"+strconv.Itoa(len(values)+1))
		values = append(values, media.Title)
	}
	if media.Description != "" {
		updateFields = append(updateFields, "description = $"+strconv.Itoa(len(values)+1))
		values = append(values, media.Description)
	}
	if media.Location != "" {
		updateFields = append(updateFields, "location = $"+strconv.Itoa(len(values)+1))
		values = append(values, media.Location)
	}
	if media.Type != "" {
		updateFields = append(updateFields, "type = $"+strconv.Itoa(len(values)+1))
		values = append(values, media.Type)
	}
	if media.MimeType != "" {
		updateFields = append(updateFields, "mimeType = $"+strconv.Itoa(len(values)+1))
		values = append(values, media.MimeType)
	}
	if media.Size != 0 {
		updateFields = append(updateFields, "size = $"+strconv.Itoa(len(values)+1))
		values = append(values, media.Size)
	}
	if media.Tags != "" {
		updateFields = append(updateFields, "tags = $"+strconv.Itoa(len(values)+1))
		values = append(values, media.Tags)
	}
	if media.MediaData != "" {
		updateFields = append(updateFields, "mediaData = $"+strconv.Itoa(len(values)+1))
		values = append(values, media.MediaData)
	}

	// Construct the full SQL query
	updateQuery := "UPDATE media SET " + strings.Join(updateFields, ", ") + " WHERE id = $" + strconv.Itoa(len(values)+1)
	values = append(values, id)

	// Execute the update
	_, err = dbConnection.Exec(context.Background(), updateQuery, values...)
	if err != nil {
		mr.logger.Error("failed to update media", slog.Any("error", err))
		return nil, fmt.Errorf("failed to update media: %w", err)
	}

	// Retrieve the updated media
	var updatedMedia models.Media
	err = dbConnection.QueryRow(context.Background(),
		"SELECT id, title, description, location, type, mimeType, size, tags, mediaData FROM media WHERE id = $1",
		id).Scan(&updatedMedia.Id, &updatedMedia.Title, &updatedMedia.Description, &updatedMedia.Location, &updatedMedia.Type, &updatedMedia.MimeType, &updatedMedia.Size, &updatedMedia.Tags, &updatedMedia.MediaData)

	if err != nil {
		mr.logger.Error("failed to retrieve updated media", slog.Any("error", err))
		return nil, fmt.Errorf("failed to retrieve updated media: %w", err)
	}

	return &updatedMedia, nil
}

func (mr *mediaRespository) Delete(id string) (bool, error) {
	dbConnection := GetDBConnection()
	if dbConnection == nil {
		mr.logger.Error("failed to get db connection")
		return false, fmt.Errorf("failed to get db connection")
	}

	// look if the media exists
	var mediaExists bool
	err := dbConnection.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM media WHERE id = $1)", id).Scan(&mediaExists)
	if err != nil {
		mr.logger.Error("failed to check if media exists", slog.Any("error", err))
		return false, fmt.Errorf("failed to check if media exists: %w", err)
	}
	if !mediaExists {
		mr.logger.Error("media with id does not exist", slog.Any("id", id))
		return false, utils.ErrMediaNotFound
	}

	_, err = dbConnection.Exec(context.Background(), "DELETE FROM media WHERE id = $1", id)
	if err != nil {
		mr.logger.Error("failed to delete media by id", slog.Any("error", err))
		return false, fmt.Errorf("failed to delete media by id: %w", err)
	}
	return true, nil
}
