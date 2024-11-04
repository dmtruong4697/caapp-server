package models

import (
	models "caapp-server/src/models/db_models"
)

type GetProfileInfoResponce struct {
	Profile models.User `json:"profile"`
}
