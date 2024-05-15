package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/procat-hq/procat-backend/internal/app/model"
)

type CartPostgres struct {
	db *sqlx.DB
}

func NewCartPostgres(db *sqlx.DB) *CartPostgres {
	return &CartPostgres{db: db}
}

func (r *CartPostgres) GetUsersCartId(userId int) (int, error) {
	query := fmt.Sprintf(`SELECT id FROM %s WHERE user_id = $1`, cartsTable)

	var id int
	err := r.db.Get(&id, query, userId)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *CartPostgres) AddItemToCart(cartId, itemId, count int) error {
	getItemCountQuery := fmt.Sprintf(`SELECT items_number FROM %s WHERE cart_id=$1 AND item_id=$2`, cartsItemsTable)

	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var itemsNumber int
	err = tx.Get(&itemsNumber, getItemCountQuery, cartId, itemId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			query := fmt.Sprintf(`INSERT INTO %s (items_number, cart_id, item_id) VALUES ($1, $2, $3)`, cartsItemsTable)
			if count == 0 {
				_, err = tx.Exec(query, 1, cartId, itemId)
			} else {
				_, err = tx.Exec(query, count, cartId, itemId)
			}
			if err != nil {
				return err
			}
			err = tx.Commit()
			return err
		}
		return err
	}

	query := fmt.Sprintf(`UPDATE %s SET items_number=$1 WHERE cart_id=$2 AND item_id=$3`, cartsItemsTable)
	if count == 0 {
		_, err = tx.Exec(query, itemsNumber+1, cartId, itemId)
	} else {
		_, err = tx.Exec(query, itemsNumber+count, cartId, itemId)
	}
	if err != nil {
		return err
	}

	err = tx.Commit()

	return err
}

func (r *CartPostgres) DeleteItemFromCart(cartId, itemId int) error {
	queryDecrease := fmt.Sprintf(`UPDATE %s SET items_number = items_number - 1
          								WHERE item_id=$1 AND cart_id=$2 RETURNING items_number`, cartsItemsTable)

	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var itemsNumber int
	err = tx.Get(&itemsNumber, queryDecrease, itemId, cartId)
	if err != nil {
		return err
	}
	if itemsNumber == 0 {
		queryToDelete := fmt.Sprintf(`DELETE FROM %s WHERE item_id=$1 AND cart_id=$2`, cartsItemsTable)
		_, err = tx.Exec(queryToDelete, itemId, cartId)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()

	return err
}

func (r *CartPostgres) GetCartItems(cartId int) ([]model.CartItem, error) {
	query := fmt.Sprintf(`SELECT i.id, i.name, i.price, c.items_number as count, COALESCE(im.image, '') AS image
									FROM %s c
									LEFT JOIN %s i ON c.item_id = i.id
									LEFT JOIN (SELECT DISTINCT ON (item_id) * FROM %s) im ON i.id = im.item_id
									WHERE c.cart_id=$1`, cartsItemsTable, itemsTable, itemsImagesTable)
	var items []model.CartItem
	err := r.db.Select(&items, query, cartId)
	if err != nil {
		return nil, err
	}
	return items, nil
}
