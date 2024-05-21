package scheduler

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/procat-hq/procat-backend/internal/app/model"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	checkInterval = 12 * time.Hour
	ordersTable   = "orders"
)

type Scheduler struct {
	db *sqlx.DB
}

func NewScheduler(db *sqlx.DB) *Scheduler {
	return &Scheduler{db: db}
}

type OrderWithDate struct {
	Id                int       `db:"id"`
	RentalPeriodStart time.Time `db:"rental_period_start"`
	RentalPeriodEnd   time.Time `db:"rental_period_end"`
}

func (s *Scheduler) CheckOrdersAndSetStatuses() {
	timer := time.NewTicker(checkInterval)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			go func() {
				query := fmt.Sprintf(`SELECT id, rental_period_start, rental_period_end
											 FROM %s WHERE status=$1 OR status=$2 OR status=$3`, ordersTable)
				var orders []OrderWithDate
				err := s.db.Select(&orders, query, model.Rent, model.Extended, model.ShouldBeReturned)
				if err != nil {
					logrus.Error("scheduler error: " + err.Error())
					return
				}

				queryChangeStatus := fmt.Sprintf(`UPDATE %s SET status=$1 WHERE id=$2`, ordersTable)
				stmt, err := s.db.Prepare(queryChangeStatus)
				if err != nil {
					logrus.Error("scheduler error: " + err.Error())
					return
				}
				defer stmt.Close()

				for _, order := range orders {
					now := time.Now()
					orderTime := time.Date(order.RentalPeriodEnd.Year(), order.RentalPeriodEnd.Month(),
						order.RentalPeriodEnd.Day(), order.RentalPeriodEnd.Hour(), order.RentalPeriodEnd.Minute(),
						order.RentalPeriodEnd.Second(), 0, time.Local)

					if now.After(orderTime) {
						_, err = stmt.Exec(model.Expired, order.Id)
						if err != nil {
							logrus.Error("scheduler error: " + err.Error())
							return
						}
					} else if now.Add(24 * time.Hour).After(orderTime) {
						_, err = stmt.Exec(model.ShouldBeReturned, order.Id)
						if err != nil {
							logrus.Error("scheduler error: " + err.Error())
							return
						}
					}
				}
				return
			}()
		}
	}
}
