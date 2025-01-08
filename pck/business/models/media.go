package models

type Media struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Type        string `json:"type"`
	MimeType    string `json:"mimeType"`
	Size        int    `json:"size"`
	Tags        string `json:"tags"`
	MediaData   string `json:"mediaData"`
}
