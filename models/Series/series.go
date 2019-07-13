package Series

import "do-mall/models"

type Series struct {
	models.Model
	BId string `json:"b_id"`
	Name string `json:"name"`
	Image string `json:"image"`
}

