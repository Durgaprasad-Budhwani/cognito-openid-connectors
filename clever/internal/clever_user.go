package internal

import "time"

// CleverUser represents a complete clever user details
type CleverUser struct {
	Created      time.Time `json:"created"`
	District     string    `json:"district"`
	Email        string    `json:"email"`
	LastModified time.Time `json:"last_modified"`
	Name         struct {
		First  string `json:"first"`
		Last   string `json:"last"`
		Middle string `json:"middle"`
	} `json:"name"`
	ID string `json:"id"`
}

type CleverUserResp struct {
	Data  CleverUser   `json:"data"`
	Links []CleverLink `json:"links"`
}
