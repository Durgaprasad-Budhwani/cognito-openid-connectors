package internal

type CleverUserInfo struct {
	Data struct {
		ID           string `json:"id"`
		District     string `json:"district"`
		Type         string `json:"type"`
		AuthorizedBy string `json:"authorized_by"`
	} `json:"data"`
	Links []struct {
		Rel string `json:"rel"`
		URI string `json:"uri"`
	} `json:"links"`
}
