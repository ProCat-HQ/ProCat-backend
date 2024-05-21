package routing

import (
	"github.com/procat-hq/procat-backend/internal/app/model"
	"math"
)

func GetRoute(response model.Api2GisResponse) (model.RoutingResult, error) {

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

	return result, nil
}

func solveTSP(distanceMatrix [][]int, durationMatrix [][]int) model.RoutingResult {
	var result model.RoutingResult
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
