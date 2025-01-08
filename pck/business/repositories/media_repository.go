package repositories

import "github.com/cosmintimis/deepfake-guardian-api/pck/business/models"

type MediaRepository interface {
	GetByID(id string) (*models.Media, error)
	Create(media *MediaPayload) (*models.Media, error)
	Update(id string, media *MediaPayload) (*models.Media, error)
	Delete(id string) (bool, error)
	GetAll() ([]models.Media, error)
}

type MediaPayload struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Type        string `json:"type"`
	MimeType    string `json:"mimeType"`
	Size        int    `json:"size"`
	Tags        string `json:"tags"`
	MediaData   string `json:"mediaData"`
}
