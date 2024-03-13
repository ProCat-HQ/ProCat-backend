package routingalgorithm

import (
	"encoding/json"
	"fmt"
	"math"
)

type Result struct {
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

type Response struct {
	GenerationTime int     `json:"generation_time"`
	Routes         []Route `json:"routes"`
}

func GetRoute(jsonInput string) (string, error) {
	var response Response
	err := json.Unmarshal([]byte(jsonInput), &response)
	if err != nil {
		return "", err
	}

	pointMap := make(map[int]bool)
	for _, route := range response.Routes {
		pointMap[route.SourceID] = true
		pointMap[route.TargetID] = true
	}
	n := len(pointMap)

	distanceMatrix := make([][]int, n)
	durationMatrix := make([][]int, n)
	for i := 0; i < n; i++ {
		distanceMatrix[i] = make([]int, n)
		durationMatrix[i] = make([]int, n)
		durationMatrix[i][i] = math.MaxInt
		distanceMatrix[i][i] = math.MaxInt
	}
	for _, route := range response.Routes {
		distanceMatrix[route.SourceID][route.TargetID] = route.Distance
		durationMatrix[route.SourceID][route.TargetID] = route.Duration
	}
	result := solveTSP(distanceMatrix, durationMatrix)
	jsonResult, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(jsonResult), nil
}

func solveTSP(distanceMatrix [][]int, durationMatrix [][]int) Result {
	var result Result
	minCost := math.MaxInt
	optimalRoute := make([]int, 0)

	visited := make([]bool, len(distanceMatrix))
	visited[0] = true // Начинаем с города 0
	currentPath := make([]int, 0, len(distanceMatrix))
	currentPath = append(currentPath, 0)

	var dfs func(int, int)
	dfs = func(currentCity, costSoFar int) {
		if len(currentPath) == len(distanceMatrix) {
			costSoFar += distanceMatrix[currentCity][0]
			if costSoFar < minCost {
				minCost = costSoFar
				optimalRoute = make([]int, len(currentPath))
				copy(optimalRoute, currentPath)
			}
			return
		}

		for nextPoint := 0; nextPoint < len(distanceMatrix); nextPoint++ {
			if !visited[nextPoint] {
				visited[nextPoint] = true
				currentPath = append(currentPath, nextPoint)
				dfs(nextPoint, costSoFar+distanceMatrix[currentCity][nextPoint])
				visited[nextPoint] = false
				currentPath = currentPath[:len(currentPath)-1]
			}
		}
	}

	dfs(0, 0)

	result.OptimalRoute = optimalRoute
	result.Distance = minCost
	duration := 0
	index := optimalRoute[0]
	for i := 0; i < len(optimalRoute)-1; i++ {
		duration += durationMatrix[optimalRoute[index]][(optimalRoute[index]+1)%len(optimalRoute)]
		index = optimalRoute[i]
	}
	result.Duration = duration
	return result
}

func main() {
	//∞  5  16 14
	//13 ∞  26 9
	//10 12 ∞  11
	//8  15 7  ∞
	jsonStr := `{
    "generation_time": 3349,
    "routes": [
        {
            "distance": 5,
            "duration": 1319,
            "source_id": 0,
            "status": "OK",
            "target_id": 1
        },
		{
            "distance": 16,
            "duration": 1319,
            "source_id": 0,
            "status": "OK",
            "target_id": 2
        },
        {
            "distance": 14,
            "duration": 603,
            "source_id": 0,
            "status": "OK",
            "target_id": 3
        },
        {
            "distance": 13,
            "duration": 1094,
            "source_id": 1,
            "status": "OK",
            "target_id": 0
        },
        {
            "distance": 26,
            "duration": 1094,
            "source_id": 1,
            "status": "OK",
            "target_id": 2
        },
		{
            "distance": 9,
            "duration": 1094,
            "source_id": 1,
            "status": "OK",
            "target_id": 3
        },
		{
            "distance": 10,
            "duration": 1094,
            "source_id": 2,
            "status": "OK",
            "target_id": 0
        },
		{
            "distance": 12,
            "duration": 1094,
            "source_id": 2,
            "status": "OK",
            "target_id": 1
        },
		{
            "distance": 11,
            "duration": 1094,
            "source_id": 2,
            "status": "OK",
            "target_id": 3
        },
		{
            "distance": 8,
            "duration": 1094,
            "source_id": 3,
            "status": "OK",
            "target_id": 0
        },
		{
            "distance": 15,
            "duration": 1094,
            "source_id": 3,
            "status": "OK",
            "target_id": 1
        },
		{
            "distance": 7,
            "duration": 1094,
            "source_id": 3,
            "status": "OK",
            "target_id": 2
        }
    ]
}`

	answer, err := GetRoute(jsonStr)
	if err != nil {
		fmt.Print("error in GetRoute")
		return
	}

	print(answer) //{"optimal_route":[0,1,3,2],"distance":31,"duration":3732}
}
