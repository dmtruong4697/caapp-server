package models

import (
	models "caapp-server/src/models/db_models"
)

type GetProfileInfoResponce struct {
	Profile models.User `json:"profile"`
}

type CheckDuplicateHashtagNameResponse struct {
	IsAvailableHashtagName bool `json:"is_available_hashtag_name"`
}
