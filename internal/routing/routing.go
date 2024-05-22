package routing

import (
	"fmt"
	"github.com/procat-hq/procat-backend/internal/app/model"
	"math"
	"sort"
)

func GetRoute(response model.Api2GisResponse, waitingHours []model.WaitingHoursForRouting) (model.RoutingResult, error) {
	if len(response.Routes) == 0 {
		return model.RoutingResult{}, fmt.Errorf("no deliveries to sort")
	}
	pointMap := make(map[int]bool)
	for _, route := range response.Routes {
		pointMap[route.SourceID] = true
		pointMap[route.TargetID] = true
	}
	n := len(pointMap)

	hours := make(map[int][]int)
	for _, hour := range waitingHours {
		if hour.Id == 0 {
			continue
		}
		if _, e := hours[hour.End]; e {
			hours[hour.End] = append(hours[hour.End], hour.Id)
		} else {
			hours[hour.End] = make([]int, 1)
			hours[hour.End][0] = hour.Id
		}
	}
	keys := make([]int, 0, len(hours))
	for key := range hours {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	distanceMatrix := make([][]int, n)
	durationMatrix := make([][]int, n)
	for i := 0; i < n; i++ {
		distanceMatrix[i] = make([]int, n)
		durationMatrix[i] = make([]int, n)
		durationMatrix[i][i] = math.MaxInt
		distanceMatrix[i][i] = math.MaxInt
	}
	for _, route := range response.Routes {
		if route.SourceID != route.TargetID {
			distanceMatrix[route.SourceID][route.TargetID] = route.Distance
			durationMatrix[route.SourceID][route.TargetID] = route.Duration
		}
	}

	var start = 0
	result := model.RoutingResult{
		OptimalRoute: []int{0},
		Distance:     0,
		Duration:     0,
	}
	for _, key := range keys {
		if len(hours[key]) == 1 {
			result.OptimalRoute = append(result.OptimalRoute, hours[key]...)
		} else {
			r := solveTSP(distanceMatrix, durationMatrix, start, hours[key])
			result.OptimalRoute = append(result.OptimalRoute, r.OptimalRoute[n-len(hours[key]):]...)
			result.Distance += r.Distance
			result.Duration += r.Duration

		}
		start = result.OptimalRoute[len(result.OptimalRoute)-1]
	}

	return result, nil
}

func solveTSP(distanceMatrix [][]int, durationMatrix [][]int, start int, points []int) model.RoutingResult {
	var result model.RoutingResult
	minCost := math.MaxInt
	optimalRoute := make([]int, 0)

	visited := make([]bool, len(distanceMatrix))
	visited[start] = true
	currentPath := make([]int, 0, len(distanceMatrix))
	currentPath = append(currentPath, start)

	f := true
	for i := 0; i < len(durationMatrix); i++ {
		for _, point := range points {
			if point == i {
				f = false
				break
			}
		}
		if f {
			for j := 0; j < len(durationMatrix); j++ {
				durationMatrix[i][j] = 0
				durationMatrix[j][i] = 0
			}
		}
		f = true
	}
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
