package repository

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/procat-hq/procat-backend/internal/app/model"
)

type SubscriptionPostgres struct {
	db *sqlx.DB
}

func NewSubscriptionPostgres(db *sqlx.DB) *SubscriptionPostgres {
	return &SubscriptionPostgres{db: db}
}

func (r *SubscriptionPostgres) GetUserSubscriptions(userId int, limit, offset int) (int, []model.Subscription, error) {
	query := fmt.Sprintf(`SELECT s.id, i.id AS item_id, i.name, i.price, i.is_in_stock, COALESCE(im.image, '') AS image
								FROM %s sub
								LEFT JOIN %s s ON sub.id=s.subscription_id
								LEFT JOIN %s i ON s.item_id = i.id
								LEFT JOIN (SELECT DISTINCT ON (item_id) * FROM %s) im ON i.id = im.item_id
								WHERE sub.user_id=$1 AND s.id IS NOT NULL
								LIMIT $2 OFFSET $3`, subscriptionsTable, subscriptionsItemsTable, itemsTable, itemsImagesTable)

	queryCount := fmt.Sprintf(`SELECT COUNT(*) FROM %s sub
										LEFT JOIN %s s ON sub.id=s.subscription_id
										WHERE sub.user_id=$1 AND s.id IS NOT NULL`,
		subscriptionsTable, subscriptionsItemsTable)

	tx, err := r.db.Beginx()
	if err != nil {
		return 0, nil, err
	}
	defer tx.Rollback()

	var subs []model.Subscription
	err = tx.Select(&subs, query, userId, limit, offset)
	if err != nil {
		return 0, nil, err
	}

	var count int
	err = tx.Get(&count, queryCount, userId)
	if err != nil {
		return 0, nil, err
	}

	return count, subs, tx.Commit()
}

func (r *SubscriptionPostgres) CreateSubscription(userId, itemId int) error {
	queryCheckItemExistence := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE id=$1`, itemsTable)
	querySub := fmt.Sprintf(`INSERT INTO %s (subscription_id, item_id)
									VALUES ((SELECT id FROM %s WHERE user_id=$1), $2)`,
		subscriptionsItemsTable, subscriptionsTable)

	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var count int
	err = tx.Get(&count, queryCheckItemExistence, itemId)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("item does not exist")
	}

	_, err = tx.Exec(querySub, userId, itemId)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *SubscriptionPostgres) DeleteSubscription(userId, subId int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1
								AND subscription_id=(SELECT id FROM %s WHERE user_id=$2)`,
		subscriptionsItemsTable, subscriptionsTable)

	_, err := r.db.Exec(query, subId, userId)
	return err
}
