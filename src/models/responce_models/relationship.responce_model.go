package models

import (
	db_models "caapp-server/src/models/db_models"
)

type GetRelationshipResponse struct {
	Status       string           `json:"status"`
	Relationship db_models.Friend `json:"relationship"`
}
