package internal

type Claim struct {
	Name       string `json:"name"`
	FistName   string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	DistrictID string `json:"district_id"`
	UserName   string `json:"username"`
}
