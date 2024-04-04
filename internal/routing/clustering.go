package routing

import (
	"github.com/muesli/clusters"
	"github.com/muesli/kmeans"
	"github.com/procat-hq/procat-backend/internal/app/model"
	"strconv"
)

func ClusterOrders(deliveries []model.DeliveryAndOrder, deliveryMen []model.DeliveryMan) (map[model.Point]int, error) {
	var coordinates clusters.Observations
	mapId := make(map[model.Point]int)
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
		mapId[model.Point{x, y, 0}] = deliveries[i].Id
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
	answer := make(map[model.Point]int)
	for i, c := range partition {
		for _, k := range c.Observations {
			answer[model.Point{Latitude: k.Coordinates()[0], Longitude: k.Coordinates()[1],
				DeliveryId: mapId[model.Point{Latitude: k.Coordinates()[0], Longitude: k.Coordinates()[1]}]}] = deliveryMen[i].Id
		}
	}

	return answer, nil
}
