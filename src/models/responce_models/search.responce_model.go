package models

type SearchUserByHashtagNameResponse struct {
	User    GetUserInfoResponce `json:"user"`
	IsFound bool                `json:"is_found"`
}
