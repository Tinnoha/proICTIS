package entity

type Enviroment struct {
	Name             string `json:"name"`
	Description      string `json:"description"`
	PhotoURL         string `json:"photo_url"`
	TypeOfEnviroment string `json:"type"`
	Auditory         string `json:"auditory"`
	IsActive         bool   `json:"is_active"`
}
