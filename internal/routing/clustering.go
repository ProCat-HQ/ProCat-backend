package routing

import (
	"fmt"
	"github.com/muesli/clusters"
	"github.com/muesli/kmeans"
	"github.com/procat-hq/procat-backend/internal/app/model"
	"strconv"
)

type Point struct {
	X float64
	Y float64
}

func ClusterOrders(deliveries []model.DeliveryAndOrder, deliveryMen []model.DeliveryMan) (map[int]int, error) {
	var coordinates clusters.Observations
	mapId := make(map[Point]int)
	for i := 0; i < len(deliveries); i++ {
		x, err := strconv.ParseFloat(deliveries[i].Latitude, 64)
		if err != nil {
			return nil, err
		}
		y, err := strconv.ParseFloat(deliveries[i].Longitude, 64)
		if err != nil {
			return nil, err
		}
		coordinates = append(coordinates, clusters.Coordinates{
			x,
			y,
		})
		mapId[Point{x, y}] = deliveries[i].Id
	}
	//km, err := kmeans.NewWithOptions(0.01, plotter.SimplePlotter{})
	//if err != nil {
	//	return nil, err
	//}
	km := kmeans.New()
	partition, err := km.Partition(coordinates, len(deliveryMen))
	if err != nil {
		return nil, err
	}
	answer := make(map[int]int)
	for i, c := range partition {
		for _, k := range c.Observations {
			answer[mapId[Point{k.Coordinates()[0], k.Coordinates()[1]}]] = deliveryMen[i].Id
		}

		fmt.Printf("Centered at x: %.6f y: %.6f\n", c.Center[0], c.Center[1])
		fmt.Printf("Matching data points: %+v\n\n", c.Observations)
	}

	return answer, nil
}
