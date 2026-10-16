package model

// Music представляет музыкальный трек
type Music struct {
	Title        string `json:"title"`
	Author       string `json:"author"`
	Duration     string `json:"duration"`
	DownloadLink string `json:"download_link"`
}
