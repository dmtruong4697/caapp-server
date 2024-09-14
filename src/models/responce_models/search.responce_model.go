package models

type SearchUserByHashtagNameResponce struct {
	Users []GetUserInfoResponce `json:"users"`
}
