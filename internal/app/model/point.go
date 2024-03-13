package model

type Point struct {
	Lat float64 `json:"lat" binding:"required"`
	Lon float64 `json:"lon" binding:"required"`
}

type RouteList struct {
	Points []Point `json:"points" binding:"required"`
}
