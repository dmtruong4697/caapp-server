package models

type GetLanguageDataListResponseItem struct {
	Code        string `json:"code"`
	EnglishName string `json:"english_name"`
	NativeName  string `json:"native_name"`
	Flag        string `json:"flag"`
}

type GetLanguageDataListResponse struct {
	Languages []GetLanguageDataListResponseItem `json:"languages"`
}
