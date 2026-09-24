package model

type City struct {
	Name string `json:"name"`
}

type CitiesResponse struct {
	Cities []City `json:"geonames"`
}
