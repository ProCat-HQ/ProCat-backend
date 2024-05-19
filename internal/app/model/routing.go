package model

type RoutingResult struct {
	OptimalRoute []int `json:"optimal_route"`
	Distance     int   `json:"distance"`
	Duration     int   `json:"duration"`
}

type Route struct {
	Distance int    `json:"distance"`
	Duration int    `json:"duration"`
	SourceID int    `json:"source_id"`
	Status   string `json:"status"`
	TargetID int    `json:"target_id"`
}

type Api2GisResponse struct {
	GenerationTime int     `json:"generation_time"`
	Routes         []Route `json:"routes"`
}
