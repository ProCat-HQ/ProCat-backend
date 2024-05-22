package routing

import (
	"github.com/muesli/clusters"
	"github.com/muesli/kmeans"
	"github.com/procat-hq/procat-backend/internal/app/model"
	"sort"
	"strconv"
)

func ClusterOrders(deliveries []model.DeliveryAndOrder, deliverymen []model.DeliveryMan) (map[model.Point]int, error) {

	if len(deliveries) == 0 {
		return nil, nil
	}
	mapId := make(map[model.Point]int)
	mapAddress := make(map[model.Point]string)
	answer := make(map[model.Point]int)
	mapOneDelivery := make(map[model.Point]int)
	mapDelivery := make(map[int]model.Point)
	mapIdCount := make(map[int]int)
	mapDeprecation := make(map[int]int)

	workStart := deliverymen[0].WorkingHoursStart.Hour()
	workEnd := deliverymen[0].WorkingHoursEnd.Hour()
	for _, man := range deliverymen {
		currentHourStart := man.WorkingHoursStart.Hour()
		currentHourEnd := man.WorkingHoursEnd.Hour()
		if currentHourEnd > workEnd {
			workEnd = currentHourEnd
		}
		if currentHourStart < workStart {
			workStart = currentHourStart
		}
		mapIdCount[man.Id] = 0
	}

	sort.Slice(deliveries, func(i, j int) bool {
		return deliveries[i].TimeStart.Before(deliveries[j].TimeStart)
	})

	for i := workStart; i < workEnd; i++ {
		var coordinates clusters.Observations
		for _, delivery := range deliveries {
			if mapDeprecation[delivery.Id] == 1 || delivery.TimeEnd.Hour() <= i {
				continue
			} else if delivery.TimeStart.Hour() > i {
				break
			}

			x, err := strconv.ParseFloat(delivery.Latitude, 64)
			if err != nil {
				return nil, err
			}
			y, err := strconv.ParseFloat(delivery.Longitude, 64)
			if err != nil {
				return nil, err
			}
			coordinates = append(coordinates, clusters.Coordinates{
				x,
				y,
			})
			mapId[model.Point{Longitude: y, Latitude: x}] = delivery.Id
			mapAddress[model.Point{Longitude: y, Latitude: x}] = delivery.Address
		}
		var deliverymenWith0 []int
		for _, man := range deliverymen {
			if man.WorkingHoursStart.Hour() <= i && man.WorkingHoursEnd.Hour() > i {
				if mapIdCount[man.Id] == 0 {
					deliverymenWith0 = append(deliverymenWith0, man.Id)
				}
			}
		}
		if len(coordinates) <= len(deliverymenWith0) {
			for i2, coordinate := range coordinates {
				answer[model.Point{
					Address:    mapAddress[model.Point{Latitude: coordinate.Coordinates()[0], Longitude: coordinate.Coordinates()[1]}],
					Latitude:   coordinate.Coordinates()[0],
					Longitude:  coordinate.Coordinates()[1],
					DeliveryId: mapId[model.Point{Latitude: coordinate.Coordinates()[0], Longitude: coordinate.Coordinates()[1]}],
				}] = deliverymenWith0[i2]

				v, e := mapIdCount[deliverymenWith0[i2]]
				if e {
					if v == 0 {
						mapIdCount[deliverymenWith0[i2]]++
						mapOneDelivery[model.Point{
							Latitude:  coordinate.Coordinates()[0],
							Longitude: coordinate.Coordinates()[1],
						}] = deliverymenWith0[i2]
						mapDelivery[deliverymenWith0[i2]] = model.Point{
							Latitude:  coordinate.Coordinates()[0],
							Longitude: coordinate.Coordinates()[1],
						}
					}
				}
				mapDeprecation[mapId[model.Point{Latitude: coordinate.Coordinates()[0], Longitude: coordinate.Coordinates()[1]}]] = 1
			}
			continue
		}
		var deliverymenForThisHour []int
		for _, man := range deliverymen {
			if man.WorkingHoursStart.Hour() <= i && man.WorkingHoursEnd.Hour() > i {
				deliverymenForThisHour = append(deliverymenForThisHour, man.Id)
				if mapIdCount[man.Id] > 0 {
					if v, e := mapDelivery[man.Id]; e {
						coordinates = append(coordinates, clusters.Coordinates{
							v.Latitude,
							v.Longitude,
						})
					}

				}
			}
		}
		if len(deliverymenForThisHour) == 0 && len(coordinates) > 0 {
			continue
		}
		km := kmeans.New()
		partition, err := km.Partition(coordinates, len(deliverymenForThisHour))
		if err != nil {
			return nil, err
		}
		for _, c := range partition {
			var deliverymanId int
			for _, k := range c.Observations {
				value, exists := mapOneDelivery[model.Point{Latitude: k.Coordinates()[0], Longitude: k.Coordinates()[1]}]
				var f = false
				if exists {
					deliverymanId = value
					for j, man := range deliverymenForThisHour {
						if man == deliverymanId {
							deliverymenForThisHour[j] = -1
							f = true
							break
						}
					}
				}
				if f {
					break
				}
			}
			if deliverymanId == 0 {
				for _, man := range deliverymenForThisHour {
					if man != -1 {
						deliverymanId = man
					}
				}
			}

			for _, k := range c.Observations {
				answer[model.Point{
					Address:    mapAddress[model.Point{Latitude: k.Coordinates()[0], Longitude: k.Coordinates()[1]}],
					Latitude:   k.Coordinates()[0],
					Longitude:  k.Coordinates()[1],
					DeliveryId: mapId[model.Point{Latitude: k.Coordinates()[0], Longitude: k.Coordinates()[1]}],
				}] = deliverymanId
				v, e := mapIdCount[deliverymanId]
				if e {
					if v == 0 {
						mapIdCount[deliverymanId]++
						mapOneDelivery[model.Point{
							Latitude:  k.Coordinates()[0],
							Longitude: k.Coordinates()[1],
						}] = deliverymanId
						mapDelivery[deliverymanId] = model.Point{
							Latitude:  k.Coordinates()[0],
							Longitude: k.Coordinates()[1],
						}
					}
				}
				mapDeprecation[mapId[model.Point{Latitude: k.Coordinates()[0], Longitude: k.Coordinates()[1]}]] = 1
			}
		}
	}
	//km, err := kmeans.NewWithOptions(0.01, plotter.SimplePlotter{})
	//if err != nil {
	//	return nil, err
	//}

	return answer, nil
}
