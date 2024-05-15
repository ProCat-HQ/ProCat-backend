package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/procat-hq/procat-backend/internal/app/model"
	"time"
)

type OrderPostgres struct {
	db *sqlx.DB
}

func NewOrderPostgres(db *sqlx.DB) *OrderPostgres {
	return &OrderPostgres{db: db}
}

func (r *OrderPostgres) GetUserById(userId int) (model.User, error) {
	query := fmt.Sprintf(`SELECT id, fullname, COALESCE(email, '') AS email, phone_number,
       							COALESCE(identification_number, '') AS identification_number, is_confirmed, role, created_at
								 FROM %s WHERE id=$1`, usersTable)

	var user model.User
	err := r.db.Get(&user, query, userId)

	return user, err
}

func (r *OrderPostgres) CreateOrder(status string, deposit bool, rpStart, rpEnd time.Time,
	address string, lat, lon float64, companyName string, userId int,
	deliveryMethod string, tStart, tEnd time.Time) (model.OrderCheque, error) {

	cartId, err := r.GetUsersCartId(userId)
	if err != nil {
		return model.OrderCheque{}, err
	}
	totalPrice, totalDeposit, err := r.GetTotalCartPrices(cartId)
	if err != nil {
		return model.OrderCheque{}, err
	}

	itemsCheque, err := r.GetItemCheque(cartId)
	if err != nil {
		return model.OrderCheque{}, err
	}

	orderCreationQuery := fmt.Sprintf(`INSERT INTO %s (status, total_price, deposit, rental_period_start, rental_period_end,
                    							address, latitude, longitude, company_name, user_id)
												VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`, ordersTable)
	depositPrice := func() int {
		if deposit {
			return totalDeposit
		}
		return 0
	}()

	tx, err := r.db.Beginx()
	if err != nil {
		return model.OrderCheque{}, err
	}
	defer tx.Rollback()

	var orderId int
	err = tx.Get(&orderId, orderCreationQuery, status, totalPrice, depositPrice, rpStart, rpEnd, address, lat, lon, companyName, userId)
	if err != nil {
		return model.OrderCheque{}, err
	}

	// if it won't work change $1 to order_id using sprintf
	queryMoveFromCartToOrder := fmt.Sprintf(`INSERT INTO %s (items_number, order_id, item_id)
													SELECT c.items_number, $1 as order_id, i.id
													FROM %s c
													LEFT JOIN %s i ON c.item_id = i.id
													WHERE c.cart_id=$2`, ordersItemsTable, cartsItemsTable, itemsTable)

	_, err = tx.Exec(queryMoveFromCartToOrder, orderId, cartId)
	if err != nil {
		return model.OrderCheque{}, err
	}

	queryDeleteFromCart := fmt.Sprintf(`DELETE FROM %s WHERE cart_id=$1`, cartsItemsTable)
	_, err = tx.Exec(queryDeleteFromCart, cartId)
	if err != nil {
		return model.OrderCheque{}, err
	}

	queryAddDelivery := fmt.Sprintf(`INSERT INTO %s (time_start, time_end, method, order_id) 
											VALUES ($1, $2, $3, $4)`, deliveriesTable)
	_, err = tx.Exec(queryAddDelivery, tStart, tEnd, deliveryMethod, orderId)
	if err != nil {
		return model.OrderCheque{}, err
	}

	queryAddPayment := fmt.Sprintf(`INSERT INTO %s (price, order_id) VALUES ($1, $2)`, paymentsTable)
	_, err = tx.Exec(queryAddPayment, totalPrice+depositPrice, orderId)

	err = tx.Commit()
	if err != nil {
		return model.OrderCheque{}, err
	}

	orderCheque := model.OrderCheque{
		OrderId:      orderId,
		TotalPrice:   totalPrice,
		TotalDeposit: depositPrice,
		Items:        itemsCheque,
	}

	return orderCheque, nil
}

func (r *OrderPostgres) GetTotalCartPrices(cartId int) (int, int, error) {
	query := fmt.Sprintf(`SELECT SUM(i.price * c.items_number) as price, SUM(i.price_deposit * c.items_number) as price_deposit FROM %s c
								LEFT JOIN %s i ON c.item_id = i.id WHERE c.cart_id=$1`, cartsItemsTable, itemsTable)
	var res struct {
		Price        int `db:"price"`
		PriceDeposit int `db:"price_deposit"`
	}
	err := r.db.Get(&res, query, cartId)
	return res.Price, res.PriceDeposit, err
}

func (r *OrderPostgres) GetUsersCartId(userId int) (int, error) {
	query := fmt.Sprintf(`SELECT id FROM %s WHERE user_id = $1`, cartsTable)

	var id int
	err := r.db.Get(&id, query, userId)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *OrderPostgres) GetItemCheque(cartId int) ([]model.ItemCheque, error) {
	query := fmt.Sprintf(`SELECT i.name, c.items_number as count, i.price, i.price_deposit FROM %s c
								LEFT JOIN %s i ON c.item_id = i.id WHERE c.cart_id=$1`, cartsItemsTable, itemsTable)
	var itemCheque []model.ItemCheque
	err := r.db.Select(&itemCheque, query, cartId)
	return itemCheque, err
}
