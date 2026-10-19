package model

<<<<<<< HEAD
// Music represents a track item with metadata and direct download link.
=======
>>>>>>> 0c59af3e9076fadbe77438343714ce617cca44d0
type Music struct {
	Title        string `json:"title"`
	Author       string `json:"author"`
	Duration     string `json:"duration"`
	DownloadLink string `json:"download_link"`
}
