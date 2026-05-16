package exercise

type CreateRequest struct {
	Name        string `json:"name"`
	MuscleGroup string `json:"muscle_group"`
	Equipment   string `json:"equipment"`
	VideoURL    string `json:"video_url"`
	Notes       string `json:"notes"`
}

type UpdateRequest struct {
	Name        string `json:"name"`
	MuscleGroup string `json:"muscle_group"`
	Equipment   string `json:"equipment"`
	VideoURL    string `json:"video_url"`
	Notes       string `json:"notes"`
}
