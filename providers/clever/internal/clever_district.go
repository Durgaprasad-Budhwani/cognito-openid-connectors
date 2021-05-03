package internal

type CleverDistrict struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type DistrictResp struct {
	District CleverDistrict `json:"data"`
	Links    []CleverLink
	URI      string
}
