package model

// Music represents a track item with metadata and direct download link.
type Music struct {
	Title        string `json:"title"`
	Author       string `json:"author"`
	Duration     string `json:"duration"`
	DownloadLink string `json:"download_link"`
	CoverURL     string `json:"cover_url,omitempty"`
}
